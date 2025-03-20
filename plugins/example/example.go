package example

import (
	"fmt"
	"gorelay/pkg/client"
	"gorelay/pkg/interfaces"
	"gorelay/pkg/packets"
	clientpackets "gorelay/pkg/packets/client"
	packetinterfaces "gorelay/pkg/packets/interfaces"
	"gorelay/pkg/packets/server"
)

// ExamplePlugin is a basic plugin that demonstrates the plugin system
type ExamplePlugin struct {
	client       *client.Client
	manager      interfaces.PluginManager
	targetName   string
	responseText string
}

// NewExamplePlugin creates a new instance of ExamplePlugin
func NewExamplePlugin() interfaces.Plugin {
	return &ExamplePlugin{
		targetName:   "Brebro",
		responseText: "Hello! This is an automated response.",
	}
}

// Initialize sets up the plugin
func (p *ExamplePlugin) Initialize(c *client.Client) error {
	p.client = c
	p.client.GetLogger().Info("HelloWorld", "Plugin initialized!")
	return nil
}

// Register registers this plugin with the plugin manager
func (p *ExamplePlugin) Register(manager interfaces.PluginManager) error {
	p.manager = manager
	manager.RegisterPlugin(p)

	// Register packet hooks
	manager.RegisterPacketHook(int32(packetinterfaces.Text), p.handleText)
	manager.RegisterPacketHook(int32(packetinterfaces.MapInfo), p.handleMapInfo)
	manager.RegisterPacketHook(int32(packetinterfaces.NewTick), p.handleNewTick)
	manager.RegisterPacketHook(int32(packetinterfaces.Update), p.handleUpdate)
	manager.RegisterPacketHook(int32(packetinterfaces.AllyShoot), p.handleAllyShoot)

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
	return "A Hello World plugin that demonstrates logging and message handling"
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

// Packet handlers
func (p *ExamplePlugin) handleText(packet packets.Packet) error {
	textPacket := packet.(*server.Text)
	p.client.GetLogger().Info("HelloWorld", "<%s> %s", textPacket.Name, textPacket.RawText)

	if textPacket.Recipient == "Extreem" {
		p.client.GetLogger().Info("HelloWorld", "Received direct message from %s: %s", textPacket.Name, textPacket.RawText)

		reply := clientpackets.NewPlayerText()
		reply.Text = fmt.Sprintf("/tell %s %s", textPacket.Name, p.responseText)

		p.client.GetLogger().Debug("HelloWorld", "Attempting to send reply packet: %s", reply.Text)

		if err := p.client.Send(reply); err != nil {
			p.client.GetLogger().Error("HelloWorld", "Failed to send reply: %v", err)
		} else {
			p.client.GetLogger().Info("HelloWorld", "Successfully sent PlayerText packet: %s", reply.Text)
		}
	}
	return nil
}

func (p *ExamplePlugin) handleMapInfo(packet packets.Packet) error {
	mapInfo := packet.(*server.MapInfo)
	p.client.GetLogger().Info("HelloWorld", "Hello from map: %s!", mapInfo.Name)
	return nil
}

func (p *ExamplePlugin) handleNewTick(packet packets.Packet) error {
	// We'll leave this empty to avoid spam
	return nil
}

func (p *ExamplePlugin) handleUpdate(packet packets.Packet) error {
	// We'll leave this empty to avoid spam
	return nil
}

func (p *ExamplePlugin) handleAllyShoot(packet packets.Packet) error {
	shoot := packet.(*server.AllyShoot)
	p.client.GetLogger().Debug("HelloWorld", "Player shot projectile: ID=%d, ContainerType=%d",
		shoot.BulletId, shoot.ContainerType)
	return nil
}

// OnUnknownPacket is called when an unknown packet is received
func (p *ExamplePlugin) OnUnknownPacket(packetID int, data []byte) {
	// Log unknown packets for debugging
	p.client.GetLogger().Debug("HelloWorld", "Received unknown packet ID %d with %d bytes",
		packetID, len(data))
}

// init registers the plugin with the plugin manager
func init() {
	// The plugin will be loaded by the plugin manager
	plugin := NewExamplePlugin()
	if plugin == nil {
		panic("failed to create example plugin")
	}
}
