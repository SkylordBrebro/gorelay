package account

import (
	"crypto"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Account represents a game account
type Account struct {
	GUID        string    `json:"guid"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Alias       string    `json:"alias"`
	ServerPref  string    `json:"serverPref"`
	CharID      int32     `json:"charId"`
	LastVerify  time.Time `json:"lastVerify"`
	Reconnect   bool      `json:"-"` // Used to signal manual reconnection
	HwidToken   string    `json:"hwidToken"`

	// Additional fields from C# implementation
	Banned                bool               `json:"banned"`
	Credits               int                `json:"credits"`
	FortuneToken          int                `json:"fortuneToken"`
	UnityCampaignPoints   int                `json:"unityCampaignPoints"`
	NextCharSlotPrice     int                `json:"nextCharSlotPrice"`
	AccountId             int64              `json:"accountId"`
	CreationTimestamp     int                `json:"creationTimestamp"`
	HasGifts              bool               `json:"hasGifts"`
	VerifiedEmail         bool               `json:"verifiedEmail"`
	PasswordError         bool               `json:"passwordError"`
	MaxNumChars           int                `json:"maxNumChars"`
	Muted                 bool               `json:"muted"`
	MutedUntil            int64              `json:"mutedUntil"`
	Originating           string             `json:"originating"`
	PetYardType           int                `json:"petYardType"`
	LastLogin             time.Time          `json:"lastLogin"`
	ForgeFireEnergy       int                `json:"forgeFireEnergy"`
	Name                  string             `json:"name"`
	NameChosen            bool               `json:"nameChosen"`
	PaymentProvider       string             `json:"paymentProvider"`
	IsAgeVerified         bool               `json:"isAgeVerified"`
	TDone                 bool               `json:"tDone"`
	SecurityQuestions     *SecurityQuestions `json:"securityQuestions"`
	AccessToken           string             `json:"accessToken"`
	AccessTokenTimestamp  int64              `json:"accessTokenTimestamp"`
	AccessTokenExpiration int                `json:"accessTokenExpiration"`
	Discoverable          bool               `json:"discoverable"`
	AccountType           string             `json:"accountType"`
	LastServer            string             `json:"lastServer"`
	Chars                 *Chars             `json:"chars"`
	DecaSignupPopup       bool               `json:"decaSignupPopup"`
	TeleportWait          int                `json:"teleportWait"`
	TOSPopup              bool               `json:"tosPopup"`
	Timestamp             string             `json:"timestamp"`
}

// SecurityQuestions represents security question settings
type SecurityQuestions struct {
	HasSecurityQuestions        bool `json:"hasSecurityQuestions" xml:"HasSecurityQuestions"`
	ShowSecurityQuestionsDialog bool `json:"showSecurityQuestionsDialog" xml:"ShowSecurityQuestionsDialog"`
}

// Chars represents character list information
type Chars struct {
	NextCharId  int    `json:"nextCharId" xml:"NextCharId"`
	MaxNumChars int    `json:"maxNumChars" xml:"MaxNumChars"`
	Characters  []Char `json:"characters" xml:"Char"`
}

// Char represents a single character
type Char struct {
	ID int `json:"id" xml:"id,attr"`
}

// TokenExpired checks if the access token has expired
func (a *Account) TokenExpired() bool {
	if a.AccessToken == "" || a.AccessToken == "0" {
		return true
	}
	return a.AccessTokenTimestamp+int64(a.AccessTokenExpiration) < time.Now().Unix()
}

// NeedAccountVerify checks if account verification is needed
func (a *Account) NeedAccountVerify() bool {
	return !a.AnyCredsError() && a.TokenExpired()
}

// NeedCharList checks if character list needs to be fetched
func (a *Account) NeedCharList() bool {
	return !a.AnyCredsError() && !a.TokenExpired() && (a.Chars == nil || a.Chars.MaxNumChars == 0)
}

// AnyCredsError checks for any credential-related errors
func (a *Account) AnyCredsError() bool {
	return a.Banned || a.PasswordError || !a.VerifiedEmail || !a.NameChosen || !a.IsAgeVerified || !a.TDone
}

// UpdateFromXML updates account information from XML response
func (a *Account) UpdateFromXML(xmlContent string) error {
	// Quick checks for known errors
	if strings.Contains(xmlContent, "passwordError") {
		a.PasswordError = true
		return nil
	}
	if strings.Contains(xmlContent, "suspended for breaching Terms of Service") {
		a.Banned = true
		return nil
	}

	// Parse XML
	var root struct {
		XMLName               xml.Name          `xml:"Account"`
		AccessToken           string            `xml:"AccessToken"`
		AccessTokenTimestamp  int64             `xml:"AccessTokenTimestamp"`
		AccessTokenExpiration int               `xml:"AccessTokenExpiration"`
		AccountID             int64             `xml:"AccountId"`
		Name                  string            `xml:"Name"`
		NameChosen            bool              `xml:"NameChosen"`
		VerifiedEmail         bool              `xml:"VerifiedEmail"`
		IsAgeVerified         bool              `xml:"IsAgeVerified"`
		SecurityQuestions     SecurityQuestions `xml:"SecurityQuestions"`
		Chars                 Chars             `xml:"Chars"`
	}

	if err := xml.Unmarshal([]byte(xmlContent), &root); err != nil {
		// If unmarshal fails, try parsing as Chars directly
		var chars Chars
		if err := xml.Unmarshal([]byte(xmlContent), &chars); err != nil {
			return fmt.Errorf("failed to parse XML response: %v", err)
		}
		a.Chars = &chars
		return nil
	}

	// Update account fields
	a.AccessToken = root.AccessToken
	a.AccessTokenTimestamp = root.AccessTokenTimestamp
	a.AccessTokenExpiration = root.AccessTokenExpiration
	a.AccountId = root.AccountID
	a.Name = root.Name
	a.NameChosen = root.NameChosen
	a.VerifiedEmail = root.VerifiedEmail
	a.IsAgeVerified = root.IsAgeVerified
	a.SecurityQuestions = &root.SecurityQuestions
	if root.Chars.MaxNumChars > 0 {
		a.Chars = &root.Chars
	}

	return nil
}

// AccountVerifyResponse represents the XML response from account verification
type AccountVerifyResponse struct {
	XMLName     xml.Name `xml:"Account"`
	AccessToken string   `xml:"AccessToken"`
}

// CharListResponse represents the XML response from char/list
type CharListResponse struct {
	XMLName xml.Name `xml:"Chars"`
	Chars   []struct {
		ID int32 `xml:"id,attr"`
	} `xml:"Char"`
	Servers struct {
		Server []struct {
			Name string `xml:"Name"`
			DNS  string `xml:"DNS"`
		} `xml:"Server"`
	} `xml:"Servers"`
}

// AccountManager handles loading and managing accounts
type AccountManager struct {
	Accounts []*Account
}

// LoadAccounts loads accounts from a JSON file
func LoadAccounts(path string) (*AccountManager, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default accounts file if it doesn't exist
			defaultAccounts := &AccountManager{
				Accounts: make([]*Account, 0),
			}
			if err := defaultAccounts.Save(path); err != nil {
				return nil, fmt.Errorf("failed to create default accounts file: %v", err)
			}
			return defaultAccounts, nil
		}
		return nil, fmt.Errorf("failed to read accounts file: %v", err)
	}

	var accounts AccountManager
	if err := json.Unmarshal(data, &accounts); err != nil {
		return nil, fmt.Errorf("failed to parse accounts file: %v", err)
	}

	// Synchronize Email and GUID fields for each account
	for _, acc := range accounts.Accounts {
		// If Email is empty but GUID is set, use GUID as Email
		if acc.Email == "" && acc.GUID != "" {
			acc.Email = acc.GUID
		}
		//todo: FIX
		log.Println("clientToken: " + string(crypto.SHA1.New().Sum([]byte(acc.Email))))
		acc.HwidToken = "b968bea6009e5d3971927d2738d329f4ea287b25"
		// If GUID is empty but Email is set, use Email as GUID
		if acc.GUID == "" && acc.Email != "" {
			acc.GUID = acc.Email
		}
	}

	return &accounts, nil
}

// Save writes the accounts to a JSON file
func (am *AccountManager) Save(path string) error {
	data, err := json.MarshalIndent(am, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal accounts: %v", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write accounts file: %v", err)
	}

	return nil
}

// AddAccount adds a new account
func (am *AccountManager) AddAccount(account *Account) {
	am.Accounts = append(am.Accounts, account)
}

// RemoveAccount removes an account by GUID
func (am *AccountManager) RemoveAccount(guid string) bool {
	for i, acc := range am.Accounts {
		if acc.GUID == guid {
			am.Accounts = append(am.Accounts[:i], am.Accounts[i+1:]...)
			return true
		}
	}
	return false
}

// GetAccount retrieves an account by GUID
func (am *AccountManager) GetAccount(guid string) *Account {
	for _, acc := range am.Accounts {
		if acc.GUID == guid {
			return acc
		}
	}
	return nil
}

// UpdateAccount updates an existing account
func (am *AccountManager) UpdateAccount(account *Account) bool {
	for i, acc := range am.Accounts {
		if acc.GUID == account.GUID {
			am.Accounts[i] = account
			return true
		}
	}
	return false
}

// VerifyAccount performs account verification and gets access token
func (a *Account) VerifyAccount(hwidToken string) error {
	// Work with a temp copy of the email for replacements/encodes
	tempEmail := a.Email
	if strings.HasPrefix(tempEmail, "steamworks") || strings.HasPrefix(tempEmail, "kongregate") {
		tempEmail = strings.ReplaceAll(tempEmail, "_", ":")
	}

	// URL encode the credentials
	encodedEmail := url.QueryEscape(tempEmail)
	encodedPassword := url.QueryEscape(a.Password)

	// Build the verification URL
	var verifyURL string
	if !strings.HasPrefix(encodedEmail, "steamworks%3A") && !strings.HasPrefix(encodedEmail, "kongregate%3A") {
		verifyURL = fmt.Sprintf("https://www.realmofthemadgod.com/account/verify?guid=%s&password=%s&clientToken=%s",
			encodedEmail, encodedPassword, hwidToken)
	} else {
		verifyURL = fmt.Sprintf("https://www.realmofthemadgod.com/account/verify?guid=%s&secret=%s&clientToken=%s",
			encodedEmail, encodedPassword, hwidToken)
	}

	// Make the request
	resp, err := http.Get(verifyURL)
	if err != nil {
		return fmt.Errorf("failed to verify account: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Check for rate limiting
	if strings.Contains(string(body), "LOGIN ATTEMPT LIMIT REACHED") {
		return fmt.Errorf("rate limit reached, please wait 5 minutes")
	}

	// Parse the XML response
	var verifyResp AccountVerifyResponse
	if err := xml.Unmarshal(body, &verifyResp); err != nil {
		return fmt.Errorf("failed to parse verify response: %v", err)
	}

	// Update account with access token
	a.AccessToken = verifyResp.AccessToken
	a.LastVerify = time.Now()

	return nil
}

// GetCharList retrieves the character list and updates account information
func (a *Account) GetCharList() error {
	if a.AccessToken == "" {
		return fmt.Errorf("no access token available")
	}

	// Build the char list URL
	charListURL := fmt.Sprintf("https://www.realmofthemadgod.com/char/list?accessToken=%s",
		url.QueryEscape(a.AccessToken))

	// Make the request
	resp, err := http.Get(charListURL)
	if err != nil {
		return fmt.Errorf("failed to get char list: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	xmlData := string(body)

	// Check for various error conditions
	if strings.Contains(xmlData, "<Error>Try again later</Error>") {
		return fmt.Errorf("rate limited, please wait 1 minute")
	}
	if strings.Contains(xmlData, "<Error>Account in use</Error>") {
		return fmt.Errorf("account currently in use")
	}

	// Parse the XML response
	var charListResp CharListResponse
	if err := xml.Unmarshal(body, &charListResp); err != nil {
		return fmt.Errorf("failed to parse char list response: %v", err)
	}

	// Update character ID if characters exist
	if len(charListResp.Chars) > 0 {
		a.CharID = charListResp.Chars[0].ID
	}

	return nil
}
