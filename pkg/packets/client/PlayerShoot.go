package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// PlayerShoot represents a packet for player shooting
type PlayerShoot struct {
	*packets.BasePacket
	Time          int32
	BulletID      byte
	ContainerType int32
	StartingPos   struct {
		X float32
		Y float32
	}
	Angle     float32
	SpeedMult float32
	LifeMult  float32
	IsBurst   bool
}

// NewPlayerShoot creates a new PlayerShoot packet
func NewPlayerShoot() *PlayerShoot {
	return &PlayerShoot{
		BasePacket: packets.NewPacket(interfaces.PlayerShoot, byte(interfaces.PlayerShoot)),
	}
}

// Type returns the packet type
func (p *PlayerShoot) Type() interfaces.PacketType {
	return interfaces.PlayerShoot
}

// Read reads the packet data from a PacketReader
func (p *PlayerShoot) Read(r *packets.PacketReader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.BulletID, err = r.ReadByte()
	if err != nil {
		return err
	}
	p.ContainerType, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.StartingPos.X, err = r.ReadFloat32()
	if err != nil {
		return err
	}
	p.StartingPos.Y, err = r.ReadFloat32()
	if err != nil {
		return err
	}
	p.Angle, err = r.ReadFloat32()
	if err != nil {
		return err
	}
	p.SpeedMult, err = r.ReadFloat32()
	if err != nil {
		return err
	}
	p.LifeMult, err = r.ReadFloat32()
	if err != nil {
		return err
	}
	p.IsBurst, err = r.ReadBool()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *PlayerShoot) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	if err := w.WriteByte(p.BulletID); err != nil {
		return err
	}
	if err := w.WriteInt32(p.ContainerType); err != nil {
		return err
	}
	if err := w.WriteFloat32(p.StartingPos.X); err != nil {
		return err
	}
	if err := w.WriteFloat32(p.StartingPos.Y); err != nil {
		return err
	}
	if err := w.WriteFloat32(p.Angle); err != nil {
		return err
	}
	if err := w.WriteFloat32(p.SpeedMult); err != nil {
		return err
	}
	if err := w.WriteFloat32(p.LifeMult); err != nil {
		return err
	}
	return w.WriteBool(p.IsBurst)
}
