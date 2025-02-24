package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// Escape represents a packet for escape actions
type Escape struct {
	*packets.BasePacket
}

// NewEscape creates a new Escape packet
func NewEscape() *Escape {
	return &Escape{
		BasePacket: packets.NewPacket(interfaces.Escape, byte(interfaces.Escape)),
	}
}

// Type returns the packet type
func (p *Escape) Type() interfaces.PacketType {
	return interfaces.Escape
}

// Read reads the packet data from a PacketReader
func (p *Escape) Read(r *packets.PacketReader) error {
	return nil
}

// Write writes the packet data to a PacketWriter
func (p *Escape) Write(w *packets.PacketWriter) error {
	return nil
}
