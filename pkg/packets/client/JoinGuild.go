package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// JoinGuild represents a packet for joining a guild
type JoinGuild struct {
	*packets.BasePacket
	GuildName string
}

// NewJoinGuild creates a new JoinGuild packet
func NewJoinGuild() *JoinGuild {
	return &JoinGuild{
		BasePacket: packets.NewPacket(interfaces.JoinGuild, byte(interfaces.JoinGuild)),
	}
}

// Type returns the packet type
func (p *JoinGuild) Type() interfaces.PacketType {
	return interfaces.JoinGuild
}

// Read reads the packet data from a PacketReader
func (p *JoinGuild) Read(r *packets.PacketReader) error {
	var err error
	p.GuildName, err = r.ReadString()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *JoinGuild) Write(w *packets.PacketWriter) error {
	return w.WriteString(p.GuildName)
}
