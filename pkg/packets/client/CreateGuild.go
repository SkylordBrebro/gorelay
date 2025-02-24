package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// CreateGuild represents a packet for creating a new guild
type CreateGuild struct {
	*packets.BasePacket
	Name string
}

// NewCreateGuild creates a new CreateGuild packet
func NewCreateGuild() *CreateGuild {
	return &CreateGuild{
		BasePacket: packets.NewPacket(interfaces.CreateGuild, byte(interfaces.CreateGuild)),
	}
}

// Type returns the packet type
func (p *CreateGuild) Type() interfaces.PacketType {
	return interfaces.CreateGuild
}

// Read reads the packet data from a PacketReader
func (p *CreateGuild) Read(r *packets.PacketReader) error {
	var err error
	p.Name, err = r.ReadString()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *CreateGuild) Write(w *packets.PacketWriter) error {
	return w.WriteString(p.Name)
}
