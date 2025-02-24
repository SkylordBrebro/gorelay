package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// GuildRemove represents a packet for removing guild members
type GuildRemove struct {
	*packets.BasePacket
	Name string
}

// NewGuildRemove creates a new GuildRemove packet
func NewGuildRemove() *GuildRemove {
	return &GuildRemove{
		BasePacket: packets.NewPacket(interfaces.GuildRemove, byte(interfaces.GuildRemove)),
	}
}

// Type returns the packet type
func (p *GuildRemove) Type() interfaces.PacketType {
	return interfaces.GuildRemove
}

// Read reads the packet data from a PacketReader
func (p *GuildRemove) Read(r *packets.PacketReader) error {
	var err error
	p.Name, err = r.ReadString()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *GuildRemove) Write(w *packets.PacketWriter) error {
	return w.WriteString(p.Name)
}
