package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// PlayerCallout represents a packet for player callouts
type PlayerCallout struct {
	*packets.BasePacket
	X float32
	Y float32
}

// NewPlayerCallout creates a new PlayerCallout packet
func NewPlayerCallout() *PlayerCallout {
	return &PlayerCallout{
		BasePacket: packets.NewPacket(interfaces.PlayerCallout, byte(interfaces.PlayerCallout)),
	}
}

// Type returns the packet type
func (p *PlayerCallout) Type() interfaces.PacketType {
	return interfaces.PlayerCallout
}

// Read reads the packet data from a PacketReader
func (p *PlayerCallout) Read(r *packets.PacketReader) error {
	var err error
	p.X, err = r.ReadFloat32()
	if err != nil {
		return err
	}
	p.Y, err = r.ReadFloat32()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *PlayerCallout) Write(w *packets.PacketWriter) error {
	if err := w.WriteFloat32(p.X); err != nil {
		return err
	}
	return w.WriteFloat32(p.Y)
}
