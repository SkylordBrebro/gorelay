package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// Load represents a packet for loading character data
type Load struct {
	*packets.BasePacket
	CharacterID  int32
	FirstSession bool
}

// NewLoad creates a new Load packet
func NewLoad() *Load {
	return &Load{
		BasePacket: packets.NewPacket(interfaces.Load, byte(interfaces.Load)),
	}
}

// Type returns the packet type
func (p *Load) Type() interfaces.PacketType {
	return interfaces.Load
}

// Read reads the packet data from a PacketReader
func (p *Load) Read(r *packets.PacketReader) error {
	var err error
	p.CharacterID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.FirstSession, err = r.ReadBool()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *Load) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.CharacterID); err != nil {
		return err
	}
	return w.WriteBool(p.FirstSession)
}
