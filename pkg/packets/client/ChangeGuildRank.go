package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// ChangeGuildRank represents a packet for changing a guild member's rank
type ChangeGuildRank struct {
	*packets.BasePacket
	Name      string
	GuildRank byte
}

// NewChangeGuildRank creates a new ChangeGuildRank packet
func NewChangeGuildRank() *ChangeGuildRank {
	return &ChangeGuildRank{
		BasePacket: packets.NewPacket(interfaces.ChangeGuildRank, byte(interfaces.ChangeGuildRank)),
	}
}

// Type returns the packet type
func (p *ChangeGuildRank) Type() interfaces.PacketType {
	return interfaces.ChangeGuildRank
}

// Read reads the packet data from a PacketReader
func (p *ChangeGuildRank) Read(r *packets.PacketReader) error {
	var err error
	p.Name, err = r.ReadString()
	if err != nil {
		return err
	}
	p.GuildRank, err = r.ReadByte()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *ChangeGuildRank) Write(w *packets.PacketWriter) error {
	if err := w.WriteString(p.Name); err != nil {
		return err
	}
	return w.WriteByte(p.GuildRank)
}
