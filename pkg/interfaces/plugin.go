package interfaces

import (
	"gorelay/pkg/client"
	"gorelay/pkg/packets"
)

// Plugin interface that must be implemented by all plugins
type Plugin interface {
	Name() string
	Initialize(client *client.Client) error
	OnEnable() error
	OnDisable() error
	Register(manager PluginManager) error
}

// PacketHook represents a packet handler function
type PacketHook func(packet packets.Packet) error

// PluginManager interface for managing plugins
type PluginManager interface {
	RegisterPlugin(plugin Plugin)
	RegisterPacketHook(packetType int32, hook PacketHook)
	UnregisterPacketHook(packetType int32, hook PacketHook)
	HandlePacket(packet packets.Packet) error
}
