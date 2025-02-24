package models

// Environment represents configuration settings that affect how the application behaves
type Environment struct {
	// Debug enables higher detail logging
	Debug bool `json:"debug"`

	// Log enables file logging
	Log bool `json:"log"`

	// LoadPlugins enables loading plugins from the default path
	LoadPlugins bool `json:"loadPlugins"`

	// LogFile specifies the path to the log file
	LogFile string `json:"logFile,omitempty"`

	// PluginPath specifies the path to load plugins from
	PluginPath string `json:"pluginPath,omitempty"`

	// ResourcePath specifies the path to game resources
	ResourcePath string `json:"resourcePath,omitempty"`

	// MaxReconnectAttempts specifies how many times to attempt reconnecting
	MaxReconnectAttempts int `json:"maxReconnectAttempts,omitempty"`

	// ReconnectDelay specifies the delay between reconnect attempts in milliseconds
	ReconnectDelay int `json:"reconnectDelay,omitempty"`
}

// DefaultEnvironment returns the default environment configuration
func DefaultEnvironment() *Environment {
	return &Environment{
		Debug:                false,
		Log:                  true,
		LoadPlugins:          true,
		LogFile:              "gorelay.log",
		PluginPath:           "plugins",
		ResourcePath:         "resources",
		MaxReconnectAttempts: 5,
		ReconnectDelay:       5000,
	}
}
