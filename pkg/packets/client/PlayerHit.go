package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// PlayerHit represents a packet for player hits
type PlayerHit struct {
	*packets.BasePacket
	BulletID int32
	ObjectID int32
}

// NewPlayerHit creates a new PlayerHit packet
func NewPlayerHit() *PlayerHit {
	return &PlayerHit{
		BasePacket: packets.NewPacket(interfaces.PlayerHit, byte(interfaces.PlayerHit)),
	}
}

// Type returns the packet type
func (p *PlayerHit) Type() interfaces.PacketType {
	return interfaces.PlayerHit
}

// Read reads the packet data from a PacketReader
func (p *PlayerHit) Read(r *packets.PacketReader) error {
	var err error
	p.BulletID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.ObjectID, err = r.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *PlayerHit) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.BulletID); err != nil {
		return err
	}
	return w.WriteInt32(p.ObjectID)
}
