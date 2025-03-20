package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gorelay/pkg/client"
	"gorelay/pkg/interfaces"
	"gorelay/pkg/models"
	"gorelay/pkg/packets"
	"gorelay/plugins/example" // Import example plugin directly
)

// PluginInstance represents a loaded plugin
type PluginInstance struct {
	Name     string
	Instance interfaces.Plugin
}

// Manager handles plugin loading and management
type Manager struct {
	plugins     []*PluginInstance
	client      *client.Client
	packetHooks map[int32][]interfaces.PacketHook
}

// NewManager creates a new plugin manager
func NewManager(client *client.Client) *Manager {
	return &Manager{
		plugins:     make([]*PluginInstance, 0),
		client:      client,
		packetHooks: make(map[int32][]interfaces.PacketHook),
	}
}

// LoadPlugin loads a plugin from the specified path
func (m *Manager) LoadPlugin(path string) error {
	m.client.GetLogger().Info("PluginManager", "Loading plugin from %s", path)

	// Get the directory containing the plugin file
	dir := filepath.Dir(path)

	// Get the package name from the directory
	pkgName := filepath.Base(dir)

	// Create a new instance of the plugin based on the package name
	var pluginInstance interfaces.Plugin
	switch pkgName {
	case "example":
		pluginInstance = example.NewExamplePlugin()
	default:
		return fmt.Errorf("unknown plugin package: %s", pkgName)
	}

	// Initialize the plugin with the client
	if err := pluginInstance.Initialize(m.client); err != nil {
		return fmt.Errorf("failed to initialize plugin: %v", err)
	}

	// Register the plugin with the manager
	if err := pluginInstance.Register(m); err != nil {
		return fmt.Errorf("failed to register plugin: %v", err)
	}

	// Enable the plugin
	if err := pluginInstance.OnEnable(); err != nil {
		return fmt.Errorf("failed to enable plugin: %v", err)
	}

	return nil
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
func (m *Manager) RegisterPacketHook(packetType int32, hook interfaces.PacketHook) {
	if hooks, exists := m.packetHooks[packetType]; exists {
		m.packetHooks[packetType] = append(hooks, hook)
	} else {
		m.packetHooks[packetType] = []interfaces.PacketHook{hook}
	}
}

// UnregisterPacketHook removes a packet handler
func (m *Manager) UnregisterPacketHook(packetType int32, hook interfaces.PacketHook) {
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

// RegisterPlugin registers a plugin with the manager
func (m *Manager) RegisterPlugin(plugin interfaces.Plugin) {
	m.plugins = append(m.plugins, &PluginInstance{
		Name:     plugin.Name(),
		Instance: plugin,
	})
	m.client.GetLogger().Info("PluginManager", "Registered plugin: %s", plugin.Name())
}

// GetPlugin returns a loaded plugin by name
func (m *Manager) GetPlugin(name string) interfaces.Plugin {
	for _, plugin := range m.plugins {
		if plugin.Name == name {
			return plugin.Instance
		}
	}
	return nil
}

// GetPlugins returns all loaded plugins
func (m *Manager) GetPlugins() []interfaces.Plugin {
	plugins := make([]interfaces.Plugin, 0, len(m.plugins))
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

// SendPacket sends a packet through the client connection
func (m *Manager) SendPacket(packet packets.Packet) error {
	if m.client == nil {
		return fmt.Errorf("client not initialized")
	}
	return m.client.Send(packet)
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
