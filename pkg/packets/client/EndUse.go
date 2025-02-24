package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// EndUse represents a packet for ending item use
type EndUse struct {
	*packets.BasePacket
	Time int32
}

// NewEndUse creates a new EndUse packet
func NewEndUse() *EndUse {
	return &EndUse{
		BasePacket: packets.NewPacket(interfaces.EndUse, byte(interfaces.EndUse)),
	}
}

// Type returns the packet type
func (p *EndUse) Type() interfaces.PacketType {
	return interfaces.EndUse
}

// Read reads the packet data from a PacketReader
func (p *EndUse) Read(r *packets.PacketReader) error {
	var err error
	p.Time, err = r.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *EndUse) Write(w *packets.PacketWriter) error {
	return w.WriteInt32(p.Time)
}
