package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// Reskin represents a packet for reskinning a character
type Reskin struct {
	*packets.BasePacket
	SkinID int32
}

// NewReskin creates a new Reskin packet
func NewReskin() *Reskin {
	return &Reskin{
		BasePacket: packets.NewPacket(interfaces.Reskin, byte(interfaces.Reskin)),
	}
}

// Type returns the packet type
func (r *Reskin) Type() interfaces.PacketType {
	return interfaces.Reskin
}

// Read reads the packet data from a PacketReader
func (r *Reskin) Read(reader *packets.PacketReader) error {
	var err error
	r.SkinID, err = reader.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (r *Reskin) Write(writer *packets.PacketWriter) error {
	return writer.WriteInt32(r.SkinID)
}
