package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// AllyShootType represents different ally shoot settings
type AllyShootType int

const (
	// HideAll hides all ally shots
	HideAll AllyShootType = iota
	// ShowAll shows all ally shots
	ShowAll
)

// ChangeAllyShoot represents a packet for changing ally shoot settings
type ChangeAllyShoot struct {
	*packets.BasePacket
	Setting int32
}

// NewChangeAllyShoot creates a new ChangeAllyShoot packet
func NewChangeAllyShoot() *ChangeAllyShoot {
	return &ChangeAllyShoot{
		BasePacket: packets.NewPacket(interfaces.ChangeAllyShoot, byte(interfaces.ChangeAllyShoot)),
	}
}

// Type returns the packet type
func (p *ChangeAllyShoot) Type() interfaces.PacketType {
	return interfaces.ChangeAllyShoot
}

// Read reads the packet data from a PacketReader
func (p *ChangeAllyShoot) Read(r *packets.PacketReader) error {
	var err error
	p.Setting, err = r.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *ChangeAllyShoot) Write(w *packets.PacketWriter) error {
	return w.WriteInt32(p.Setting)
}
