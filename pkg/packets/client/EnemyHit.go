package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// EnemyHit represents a packet for enemy hit events
type EnemyHit struct {
	*packets.BasePacket
	Time     int32
	BulletID uint16
	SourceID int32
	TargetID int32
	Killed   bool
	OwnerID  int32
}

// NewEnemyHit creates a new EnemyHit packet
func NewEnemyHit() *EnemyHit {
	return &EnemyHit{
		BasePacket: packets.NewPacket(interfaces.EnemyHit, byte(interfaces.EnemyHit)),
	}
}

// Type returns the packet type
func (p *EnemyHit) Type() interfaces.PacketType {
	return interfaces.EnemyHit
}

// Read reads the packet data from a PacketReader
func (p *EnemyHit) Read(r *packets.PacketReader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.BulletID, err = r.ReadUInt16()
	if err != nil {
		return err
	}
	p.SourceID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.TargetID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Killed, err = r.ReadBool()
	if err != nil {
		return err
	}
	p.OwnerID, err = r.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *EnemyHit) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	if err := w.WriteUInt16(p.BulletID); err != nil {
		return err
	}
	if err := w.WriteInt32(p.SourceID); err != nil {
		return err
	}
	if err := w.WriteInt32(p.TargetID); err != nil {
		return err
	}
	if err := w.WriteBool(p.Killed); err != nil {
		return err
	}
	return w.WriteInt32(p.OwnerID)
}
