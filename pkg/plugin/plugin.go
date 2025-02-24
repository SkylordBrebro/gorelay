package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"gorelay/pkg/client"
	"gorelay/pkg/models"
	"gorelay/pkg/packets"
)

// Plugin interface that must be implemented by all plugins
type Plugin interface {
	Name() string
	Initialize(client *client.Client) error
	OnEnable() error
	OnDisable() error
}

// PacketHook represents a packet handler function
type PacketHook func(packet packets.Packet) error

// PluginInstance represents a loaded plugin
type PluginInstance struct {
	Name     string
	Instance Plugin
	client   *client.Client
}

// Manager handles plugin loading and management
type Manager struct {
	plugins     []*PluginInstance
	client      *client.Client
	packetHooks map[int32][]PacketHook
}

// NewManager creates a new plugin manager
func NewManager(client *client.Client) *Manager {
	return &Manager{
		plugins:     make([]*PluginInstance, 0),
		client:      client,
		packetHooks: make(map[int32][]PacketHook),
	}
}

// LoadPlugin loads a plugin from the specified path
func (m *Manager) LoadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %v", err)
	}

	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin does not export 'Plugin': %v", err)
	}

	instance, ok := symPlugin.(Plugin)
	if !ok {
		return fmt.Errorf("invalid plugin type")
	}

	// Initialize the plugin
	if err := instance.Initialize(m.client); err != nil {
		return fmt.Errorf("failed to initialize plugin: %v", err)
	}

	plugin := &PluginInstance{
		Name:     instance.Name(),
		Instance: instance,
		client:   m.client,
	}

	m.plugins = append(m.plugins, plugin)

	// Enable the plugin
	return instance.OnEnable()
}

// UnloadPlugin disables and unloads a plugin
func (m *Manager) UnloadPlugin(name string) error {
	for i, plugin := range m.plugins {
		if plugin.Name == name {
			if err := plugin.Instance.OnDisable(); err != nil {
				return err
			}

			// Remove the plugin from the slice
			m.plugins = append(m.plugins[:i], m.plugins[i+1:]...)
			return nil
		}
	}
	return nil
}

// RegisterPacketHook registers a packet handler for a specific packet type
func (m *Manager) RegisterPacketHook(packetType int32, hook PacketHook) {
	if hooks, exists := m.packetHooks[packetType]; exists {
		m.packetHooks[packetType] = append(hooks, hook)
	} else {
		m.packetHooks[packetType] = []PacketHook{hook}
	}
}

// UnregisterPacketHook removes a packet handler
func (m *Manager) UnregisterPacketHook(packetType int32, hook PacketHook) {
	if hooks, exists := m.packetHooks[packetType]; exists {
		for i, h := range hooks {
			if &h == &hook {
				m.packetHooks[packetType] = append(hooks[:i], hooks[i+1:]...)
				break
			}
		}
	}
}

// HandlePacket processes a packet through all registered hooks
func (m *Manager) HandlePacket(packet packets.Packet) error {
	if hooks, exists := m.packetHooks[int32(packet.ID())]; exists {
		for _, hook := range hooks {
			if err := hook(packet); err != nil {
				return err
			}
		}
	}
	return nil
}

// SwitchServer switches the client to the specified server
func (m *Manager) SwitchServer(serverName string) error {
	return m.client.SwitchServer(serverName)
}

// GetCurrentServer returns the current server configuration
func (m *Manager) GetCurrentServer() *models.Server {
	return m.client.GetCurrentServer()
}

// GetAvailableServers returns a list of all available servers
func (m *Manager) GetAvailableServers() models.ServerList {
	if models.CachedServers != nil {
		return models.CachedServers
	}
	return models.ServerList{models.DefaultServer.Name: models.DefaultServer}
}

// GetPlugin returns a loaded plugin by name
func (m *Manager) GetPlugin(name string) Plugin {
	for _, plugin := range m.plugins {
		if plugin.Name == name {
			return plugin.Instance
		}
	}
	return nil
}

// GetPlugins returns all loaded plugins
func (m *Manager) GetPlugins() []Plugin {
	plugins := make([]Plugin, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		plugins = append(plugins, plugin.Instance)
	}
	return plugins
}

// LoadPluginsFromDirectory loads all enabled plugins from the specified directory
func (m *Manager) LoadPluginsFromDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %v", err)
	}

	for _, entry := range entries {
		// Skip directories and non-Go files
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		// Skip disabled plugins (prefixed with -)
		if strings.HasPrefix(entry.Name(), "-") {
			m.client.GetLogger().Info("PluginManager", "Skipping disabled plugin: %s", entry.Name())
			continue
		}

		pluginPath := filepath.Join(dir, entry.Name())
		if err := m.LoadPlugin(pluginPath); err != nil {
			m.client.GetLogger().Error("PluginManager", "Failed to load plugin %s: %v", entry.Name(), err)
			continue
		}
	}

	return nil
}
