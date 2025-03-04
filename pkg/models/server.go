package models

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Server represents a game server that can be connected to
type Server struct {
	Name     string  `json:"name" xml:"name"`
	Address  string  `json:"address"`
	DNS      string  `json:"dns,omitempty" xml:"dns"`
	Port     int     `json:"port"`
	Usage    float32 `json:"usage" xml:"usage"`
	MaxUsers int32   `json:"maxUsers"`
	Lat      float32 `json:"lat,omitempty" xml:"lat"`
	Long     float32 `json:"long,omitempty" xml:"long"`
}

// ServerList is a map of server names to their configurations
type ServerList map[string]*Server

// XMLServerList represents the XML response from the server list API
type XMLServerList struct {
	XMLName xml.Name    `xml:"Servers"`
	Servers []XMLServer `xml:"Server"`
}

type XMLServer struct {
	Name  string  `xml:"Name"`
	DNS   string  `xml:"DNS"`
	Lat   float32 `xml:"Lat"`
	Long  float32 `xml:"Long"`
	Usage float32 `xml:"Usage"`
}

// DefaultServer is the fallback server if no others are available
var DefaultServer = &Server{
	Name:    "USEast",
	Address: "54.234.226.24", // Keeping one default IP as absolute fallback
	Port:    2050,
}

// CachedServers stores the last fetched server list
var CachedServers ServerList

// FetchServers retrieves the current server list from the ROTMG API
func FetchServers(guid string, password string) (ServerList, error) {
	// Check for empty credentials
	if guid == "" {
		return nil, fmt.Errorf("empty email/guid provided")
	}
	if password == "" {
		return nil, fmt.Errorf("empty password provided")
	}

	// URL encode the guid and password parameters
	encodedGuid := url.QueryEscape(guid)
	encodedPassword := url.QueryEscape(password)

	requestURL := fmt.Sprintf("https://www.realmofthemadgod.com/account/servers?guid=%s&password=%s",
		encodedGuid, encodedPassword)

	// Log the request URL for debugging (without the password)
	fmt.Printf("Fetching servers from URL: https://www.realmofthemadgod.com/account/servers?guid=%s&password=REDACTED\n", encodedGuid)

	// Create a custom HTTP client with appropriate headers
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add user agent and other headers that might be required
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/xml, text/xml, */*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch servers: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Log the response body for debugging
	fmt.Printf("Got server list")
	//fmt.Printf("Response status: %s\n", resp.Status)
	//fmt.Printf("Response body: %s\n", string(body))

	// Check for common error responses
	responseStr := string(body)
	if responseStr == "<Error>Error.incorrectEmailOrPassword</Error>" {
		return nil, fmt.Errorf("authentication failed: incorrect email/guid or password")
	}

	// Parse XML response
	var xmlList XMLServerList
	if err := xml.Unmarshal(body, &xmlList); err != nil {
		// Try to parse as error response
		var errorResp struct {
			XMLName xml.Name `xml:"Error"`
			Message string   `xml:",chardata"`
		}
		if xmlErr := xml.Unmarshal(body, &errorResp); xmlErr == nil {
			return nil, fmt.Errorf("server returned error: %s", errorResp.Message)
		}
		return nil, fmt.Errorf("failed to parse server list: %v", err)
	}

	// Convert to ServerList format
	servers := make(ServerList)
	for _, s := range xmlList.Servers {
		servers[s.Name] = &Server{
			Name:    s.Name,
			Address: s.DNS, // Use DNS as the address
			DNS:     s.DNS,
			Port:    2050, // Default ROTMG port
			Usage:   s.Usage,
			Lat:     s.Lat,
			Long:    s.Long,
		}
	}

	// Cache the servers for future use
	CachedServers = servers
	return servers, nil
}

// GetServer returns a server configuration by name
func GetServer(name string) *Server {
	if CachedServers != nil {
		if server, ok := CachedServers[name]; ok {
			return server
		}
		// If server name not found but we have cached servers, return first available
		for _, server := range CachedServers {
			return server
		}
	}
	return DefaultServer // Absolute fallback if no servers available
}

// ServerList represents a list of available game servers
type ServerListStruct struct {
	Servers  []Server `json:"servers"`
	MaxUsers int32    `json:"maxUsers"`
}

// ServerStats represents server statistics
type ServerStats struct {
	Name           string
	Address        string
	ConnectedUsers int32
	Uptime         int64
	LastUpdate     int64
	CPU            float32
	Memory         float32
	Bandwidth      float32
}

// ServerConfig represents server configuration
type ServerConfig struct {
	// Local server settings
	LocalEnabled bool   `json:"localEnabled"`
	LocalPort    int    `json:"localPort"`
	LocalHost    string `json:"localHost"`

	// Remote server settings
	RemoteEnabled bool   `json:"remoteEnabled"`
	RemoteHost    string `json:"remoteHost"`
	RemotePort    int    `json:"remotePort"`

	// Server behavior
	MaxConnections int32 `json:"maxConnections"`
	Timeout        int32 `json:"timeout"`
	KeepAlive      bool  `json:"keepAlive"`
	Debug          bool  `json:"debug"`
}
