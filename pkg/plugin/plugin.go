package plugin

import (
	"fmt"
	"plugin"
	"reflect"

	"gorelay/pkg/client"
	"gorelay/pkg/packets"
)

// Plugin defines the interface that all plugins must implement
type Plugin interface {
	// Initialize is called when the plugin is first loaded
	Initialize(client *client.Client) error

	// Name returns the name of the plugin
	Name() string

	// Author returns the author of the plugin
	Author() string

	// Version returns the plugin version
	Version() string

	// Description returns the plugin description
	Description() string

	// OnEnable is called when the plugin is enabled
	OnEnable() error

	// OnDisable is called when the plugin is disabled
	OnDisable() error
}

// PacketHook represents a method that handles a specific packet type
type PacketHook struct {
	Plugin     Plugin
	Method     reflect.Method
	PacketType reflect.Type
}

// Manager handles plugin loading and lifecycle
type Manager struct {
	plugins     map[string]Plugin
	client      *client.Client
	packetHooks map[int32][]PacketHook
}

// NewManager creates a new plugin manager
func NewManager(client *client.Client) *Manager {
	return &Manager{
		plugins:     make(map[string]Plugin),
		client:      client,
		packetHooks: make(map[int32][]PacketHook),
	}
}

// LoadPlugin loads and initializes a plugin
func (m *Manager) LoadPlugin(path string) error {
	// Load the plugin
	plug, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to load plugin: %v", err)
	}

	// Look up the plugin symbol
	symPlugin, err := plug.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin does not export 'Plugin' symbol: %v", err)
	}

	// Assert that the symbol is a Plugin
	p, ok := symPlugin.(Plugin)
	if !ok {
		return fmt.Errorf("symbol is not a Plugin")
	}

	// Initialize the plugin
	if err := p.Initialize(m.client); err != nil {
		return fmt.Errorf("failed to initialize plugin: %v", err)
	}

	// Register packet hooks
	if err := m.registerHooks(p); err != nil {
		return fmt.Errorf("failed to register hooks: %v", err)
	}

	m.plugins[p.Name()] = p
	return p.OnEnable()
}

// UnloadPlugin disables and unloads a plugin
func (m *Manager) UnloadPlugin(name string) error {
	p, exists := m.plugins[name]
	if !exists {
		return nil
	}

	if err := p.OnDisable(); err != nil {
		return err
	}

	// Unregister packet hooks
	m.unregisterHooks(p)

	delete(m.plugins, name)
	return nil
}

// registerHooks registers all packet hooks for a plugin
func (m *Manager) registerHooks(p Plugin) error {
	pType := reflect.TypeOf(p)

	// Iterate through all methods
	for i := 0; i < pType.NumMethod(); i++ {
		method := pType.Method(i)

		// Check if method is a packet hook
		if hook := m.parseHook(p, method); hook != nil {
			packetID := m.getPacketID(hook.PacketType)
			m.packetHooks[packetID] = append(m.packetHooks[packetID], *hook)
		}
	}

	return nil
}

// unregisterHooks removes all packet hooks for a plugin
func (m *Manager) unregisterHooks(p Plugin) {
	for id, hooks := range m.packetHooks {
		filtered := make([]PacketHook, 0)
		for _, hook := range hooks {
			if hook.Plugin != p {
				filtered = append(filtered, hook)
			}
		}
		if len(filtered) > 0 {
			m.packetHooks[id] = filtered
		} else {
			delete(m.packetHooks, id)
		}
	}
}

// parseHook checks if a method is a packet hook and returns its info
func (m *Manager) parseHook(p Plugin, method reflect.Method) *PacketHook {
	mType := method.Type

	// Must have exactly one parameter (besides receiver)
	if mType.NumIn() != 2 {
		return nil
	}

	// Parameter must implement packets.Packet
	paramType := mType.In(1)
	if !paramType.Implements(reflect.TypeOf((*packets.Packet)(nil)).Elem()) {
		return nil
	}

	return &PacketHook{
		Plugin:     p,
		Method:     method,
		PacketType: paramType,
	}
}

// getPacketID gets the packet ID from a packet type
func (m *Manager) getPacketID(t reflect.Type) int32 {
	// Create a new instance of the packet type
	packet := reflect.New(t.Elem()).Interface().(packets.Packet)
	return packet.ID()
}

// HandlePacket dispatches a packet to all registered hooks
func (m *Manager) HandlePacket(packet packets.Packet) {
	hooks, exists := m.packetHooks[packet.ID()]
	if !exists {
		return
	}

	for _, hook := range hooks {
		// Call the hook method
		hook.Method.Func.Call([]reflect.Value{
			reflect.ValueOf(hook.Plugin),
			reflect.ValueOf(packet),
		})
	}
}

// GetPlugin returns a loaded plugin by name
func (m *Manager) GetPlugin(name string) Plugin {
	return m.plugins[name]
}

// GetPlugins returns all loaded plugins
func (m *Manager) GetPlugins() []Plugin {
	plugins := make([]Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p)
	}
	return plugins
}
