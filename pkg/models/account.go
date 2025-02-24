package models

// AccountInfo represents the account configuration used at startup
type AccountInfo struct {
	BuildVersion string             `json:"buildVersion"`
	LocalServer  *LocalServerConfig `json:"localServer,omitempty"`
	Accounts     []Account          `json:"accounts"`
}

// Account represents a game account with its credentials and settings
type Account struct {
	Alias      string         `json:"alias"`
	Email      string         `json:"email"`
	GUID       string         `json:"guid"`
	Password   string         `json:"password"`
	ServerPref string         `json:"serverPref"`
	CharInfo   *CharacterInfo `json:"charInfo,omitempty"`
	Proxy      *ProxyConfig   `json:"proxy,omitempty"`
	Pathfinder bool           `json:"pathfinder,omitempty"`
	Reconnect  bool           `json:"-"` // Used to signal manual reconnection
}

// CharacterInfo contains information about an account's characters
type CharacterInfo struct {
	CharID      int32 `json:"charId"`
	NextCharID  int32 `json:"nextCharId"`
	MaxNumChars int32 `json:"maxNumChars"`
}

// LocalServerConfig contains settings for the local proxy server
type LocalServerConfig struct {
	Enabled bool `json:"enabled"`
	Port    int  `json:"port,omitempty"`
}

// ProxyConfig contains proxy server settings
type ProxyConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// AccountInUseError represents an error when an account is already in use
type AccountInUseError struct {
	Account *Account
	Server  string
}

// Error implements the error interface
func (e *AccountInUseError) Error() string {
	return "Account " + e.Account.Alias + " is already in use on server " + e.Server
}
