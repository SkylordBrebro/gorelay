package models

// Server represents a game server that can be connected to
type Server struct {
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	Port     int     `json:"port"`
	Usage    int32   `json:"usage"`
	MaxUsers int32   `json:"maxUsers"`
	Lat      float32 `json:"lat,omitempty"`
	Long     float32 `json:"long,omitempty"`
	DNS      string  `json:"dns,omitempty"`
}

// ServerList is a map of server names to their configurations
type ServerList map[string]*Server

// DefaultServers contains the list of available game servers
var DefaultServers = ServerList{
	"EUEast":      {Name: "EUEast", Address: "18.184.218.174"},
	"EUSouthWest": {Name: "EUSouthWest", Address: "35.180.67.120"},
	"EUNorth":     {Name: "EUNorth", Address: "18.159.133.120"},
	"USWest4":     {Name: "USWest4", Address: "54.235.235.140"},
	"EUWest2":     {Name: "EUWest2", Address: "52.16.86.215"},
	"USSouth3":    {Name: "USSouth3", Address: "52.207.206.31"},
	"Asia":        {Name: "Asia", Address: "3.0.147.127"},
	"EUWest":      {Name: "EUWest", Address: "15.237.60.223"},
	"USMidWest":   {Name: "USMidWest", Address: "18.221.120.59"},
	"USSouth":     {Name: "USSouth", Address: "3.82.126.16"},
	"USWest3":     {Name: "USWest3", Address: "18.144.30.153"},
	"USSouthWest": {Name: "USSouthWest", Address: "54.153.13.68"},
	"Australia":   {Name: "Australia", Address: "54.79.72.84"},
	"USWest":      {Name: "USWest", Address: "54.86.47.176"},
	"USMidWest2":  {Name: "USMidWest2", Address: "3.140.254.133"},
	"USNorthWest": {Name: "USNorthWest", Address: "34.238.176.119"},
	"USEast2":     {Name: "USEast2", Address: "54.209.152.223"},
	"USEast":      {Name: "USEast", Address: "54.234.226.24"},
}

// GetServer returns a server configuration by name
func GetServer(name string) *Server {
	if server, ok := DefaultServers[name]; ok {
		return server
	}
	return DefaultServers["USEast"] // Default to USEast if server not found
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
