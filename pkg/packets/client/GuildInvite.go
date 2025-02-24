package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// GuildInvite represents a packet for guild invitations
type GuildInvite struct {
	*packets.BasePacket
	Name string
}

// NewGuildInvite creates a new GuildInvite packet
func NewGuildInvite() *GuildInvite {
	return &GuildInvite{
		BasePacket: packets.NewPacket(interfaces.GuildInvite, byte(interfaces.GuildInvite)),
	}
}

// Type returns the packet type
func (p *GuildInvite) Type() interfaces.PacketType {
	return interfaces.GuildInvite
}

// Read reads the packet data from a PacketReader
func (p *GuildInvite) Read(r *packets.PacketReader) error {
	var err error
	p.Name, err = r.ReadString()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *GuildInvite) Write(w *packets.PacketWriter) error {
	return w.WriteString(p.Name)
}
