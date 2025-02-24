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

// ServerList represents a list of available game servers
type ServerList struct {
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
