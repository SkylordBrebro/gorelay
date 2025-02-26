package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config represents the application configuration
type Config struct {
	BuildHash string `json:"buildHash"` //hex hash
	BuildVersion string `json:"buildVersion"` //version string
	LocalServer  struct {
		Enabled bool `json:"enabled"`
		Port    int  `json:"port"`
	} `json:"localServer"`
	Debug bool `json:"debug"`

	// Game settings
	AutoNexusThreshold float32 `json:"autoNexusThreshold"`
	AutoHealThreshold  float32 `json:"autoHealThreshold"`
	AutoHealMP         float32 `json:"autoHealMP"`
	ReconnectDelay     int     `json:"reconnectDelay"`
	SafeWalk           bool    `json:"safeWalk"`
	AutoAim            bool    `json:"autoAim"`

	// Proxy settings
	Proxy struct {
		Enabled  bool   `json:"enabled"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"proxy"`

	// Plugin settings
	Plugins struct {
		Enabled bool     `json:"enabled"`
		Path    string   `json:"path"`
		List    []string `json:"list"`
	} `json:"plugins"`

	HWIDToken string `json:"hwidToken"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config if it doesn't exist
			defaultConfig := &Config{
				BuildVersion: "", // Build version will be fetched dynamically
				Debug:        false,
				LocalServer: struct {
					Enabled bool `json:"enabled"`
					Port    int  `json:"port"`
				}{
					Enabled: false,
					Port:    2050,
				},
				AutoNexusThreshold: 0.3,
				AutoHealThreshold:  0.6,
				AutoHealMP:         0.4,
				ReconnectDelay:     5000,
				SafeWalk:           true,
				AutoAim:            true,
				Proxy: struct {
					Enabled  bool   `json:"enabled"`
					Host     string `json:"host"`
					Port     int    `json:"port"`
					Username string `json:"username"`
					Password string `json:"password"`
				}{
					Enabled: false,
				},
				Plugins: struct {
					Enabled bool     `json:"enabled"`
					Path    string   `json:"path"`
					List    []string `json:"list"`
				}{
					Enabled: false,
					Path:    "plugins",
					List:    make([]string, 0),
				},
			}
			if err := defaultConfig.Save(path); err != nil {
				return nil, fmt.Errorf("failed to create default config: %v", err)
			}
			return defaultConfig, nil
		}
		return nil, fmt.Errorf("failed to read config: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &config, nil
}

// Save writes the configuration to a JSON file
func (c *Config) Save(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	return nil
}

// SaveConfig saves the configuration back to the JSON file
func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
