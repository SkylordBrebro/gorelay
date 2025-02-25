package main

import (
	"gorelay/pkg/client"
	"gorelay/pkg/packets/server"
)

// ExamplePlugin is a basic plugin that demonstrates the plugin system
type ExamplePlugin struct {
	client *client.Client
}

// Plugin is the exported symbol that the plugin manager looks for
var Plugin ExamplePlugin

// Initialize sets up the plugin
func (p *ExamplePlugin) Initialize(c *client.Client) error {
	p.client = c
	p.client.GetLogger().Info("HelloWorld", "Plugin initialized!")
	return nil
}

// Name returns the plugin name
func (p *ExamplePlugin) Name() string {
	return "HelloWorld"
}

// Author returns the plugin author
func (p *ExamplePlugin) Author() string {
	return "Example Author"
}

// Version returns the plugin version
func (p *ExamplePlugin) Version() string {
	return "1.0.0"
}

// Description returns the plugin description
func (p *ExamplePlugin) Description() string {
	return "A Hello World plugin that demonstrates logging"
}

// OnEnable is called when the plugin is enabled
func (p *ExamplePlugin) OnEnable() error {
	p.client.GetLogger().Info("HelloWorld", "Hello, World! Plugin enabled and ready to go!")
	return nil
}

// OnDisable is called when the plugin is disabled
func (p *ExamplePlugin) OnDisable() error {
	p.client.GetLogger().Info("HelloWorld", "Goodbye, World! Plugin shutting down...")
	return nil
}

// OnMapInfo is called when a MapInfo packet is received
func (p *ExamplePlugin) OnMapInfo(packet *server.MapInfo) {
	p.client.GetLogger().Info("HelloWorld", "Hello from map: %s!", packet.Name)
}

// OnNewTick is called when a NewTick packet is received
func (p *ExamplePlugin) OnNewTick(packet *server.NewTick) {
	// We'll leave this empty to avoid spam
}

// OnUpdate is called when an Update packet is received
func (p *ExamplePlugin) OnUpdate(packet *server.Update) {
	// We'll leave this empty to avoid spam
}
