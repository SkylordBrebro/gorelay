package models

// PluginInfo represents information about a plugin
type PluginInfo struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled"`
}

// LibraryInfo represents information about a library
type LibraryInfo struct {
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled"`
}

// PluginConfig represents plugin configuration
type PluginConfig struct {
	Path    string   `json:"path"`
	Enabled bool     `json:"enabled"`
	List    []string `json:"list"`
}

// PluginMetadata represents metadata about a loaded plugin
type PluginMetadata struct {
	Info      PluginInfo
	Path      string
	LoadTime  int64
	Instances int32
	Events    []string
	Methods   []string
}
