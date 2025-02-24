package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// Retitle represents a packet for changing a title
type Retitle struct {
	*packets.BasePacket
	Prefix int32
	Suffix int32
}

// NewRetitle creates a new Retitle packet
func NewRetitle() *Retitle {
	return &Retitle{
		BasePacket: packets.NewPacket(interfaces.Retitle, byte(interfaces.Retitle)),
	}
}

// Type returns the packet type
func (r *Retitle) Type() interfaces.PacketType {
	return interfaces.Retitle
}

// Read reads the packet data from a PacketReader
func (r *Retitle) Read(reader *packets.PacketReader) error {
	var err error
	r.Prefix, err = reader.ReadInt32()
	if err != nil {
		return err
	}
	r.Suffix, err = reader.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (r *Retitle) Write(writer *packets.PacketWriter) error {
	if err := writer.WriteInt32(r.Prefix); err != nil {
		return err
	}
	return writer.WriteInt32(r.Suffix)
}
