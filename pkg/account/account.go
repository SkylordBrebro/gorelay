package account

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Account represents a game account
type Account struct {
	GUID       string `json:"guid"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Alias      string `json:"alias"`
	ServerPref string `json:"serverPref"`
	CharID     int32  `json:"charId"`
	Reconnect  bool   `json:"-"` // Used to signal manual reconnection
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
