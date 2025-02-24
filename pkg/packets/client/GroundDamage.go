package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// GroundDamage represents a packet for ground damage events
type GroundDamage struct {
	*packets.BasePacket
	Time     int32
	Position *Location
}

// NewGroundDamage creates a new GroundDamage packet
func NewGroundDamage() *GroundDamage {
	return &GroundDamage{
		BasePacket: packets.NewPacket(interfaces.GroundDamage, byte(interfaces.GroundDamage)),
	}
}

// Type returns the packet type
func (p *GroundDamage) Type() interfaces.PacketType {
	return interfaces.GroundDamage
}

// Read reads the packet data from a PacketReader
func (p *GroundDamage) Read(r *packets.PacketReader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Position, err = NewLocation(r)
	return err
}

// Write writes the packet data to a PacketWriter
func (p *GroundDamage) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	return p.Position.Write(w)
}
