package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// Create represents a packet for creating a new character
type Create struct {
	*packets.BasePacket
	ClassType    uint16
	SkinType     uint16
	IsChallenger bool
	IsSeasonal   bool
}

// NewCreate creates a new Create packet
func NewCreate() *Create {
	return &Create{
		BasePacket: packets.NewPacket(interfaces.Create, byte(interfaces.Create)),
	}
}

// Type returns the packet type
func (p *Create) Type() interfaces.PacketType {
	return interfaces.Create
}

// Read reads the packet data from a PacketReader
func (p *Create) Read(r *packets.PacketReader) error {
	var err error
	p.ClassType, err = r.ReadUInt16()
	if err != nil {
		return err
	}
	p.SkinType, err = r.ReadUInt16()
	if err != nil {
		return err
	}
	p.IsChallenger, err = r.ReadBool()
	if err != nil {
		return err
	}
	p.IsSeasonal, err = r.ReadBool()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *Create) Write(w *packets.PacketWriter) error {
	if err := w.WriteUInt16(p.ClassType); err != nil {
		return err
	}
	if err := w.WriteUInt16(p.SkinType); err != nil {
		return err
	}
	if err := w.WriteBool(p.IsChallenger); err != nil {
		return err
	}
	return w.WriteBool(p.IsSeasonal)
}
