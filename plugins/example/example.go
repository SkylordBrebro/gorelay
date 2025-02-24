package main

import (
	"gorelay/pkg/client"
	"gorelay/pkg/packets"
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
	return nil
}

// Name returns the plugin name
func (p *ExamplePlugin) Name() string {
	return "Example"
}

// Author returns the plugin author
func (p *ExamplePlugin) Author() string {
	return "GoRelay Team"
}

// Version returns the plugin version
func (p *ExamplePlugin) Version() string {
	return "1.0.0"
}

// Description returns the plugin description
func (p *ExamplePlugin) Description() string {
	return "An example plugin demonstrating the plugin system"
}

// OnEnable is called when the plugin is enabled
func (p *ExamplePlugin) OnEnable() error {
	return nil
}

// OnDisable is called when the plugin is disabled
func (p *ExamplePlugin) OnDisable() error {
	return nil
}

// OnMapInfo is called when a MapInfo packet is received
func (p *ExamplePlugin) OnMapInfo(packet *packets.MapInfoPacket) {
	p.client.GetLogger().Info("Example", "Entered map: %s", packet.Name)
}

// OnNewTick is called when a NewTick packet is received
func (p *ExamplePlugin) OnNewTick(packet *packets.NewTickPacket) {
	// Handle game tick updates
}

// OnUpdate is called when an Update packet is received
func (p *ExamplePlugin) OnUpdate(packet *packets.UpdatePacket) {
	// Handle entity updates
}
