package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Damage represents the server packet for damage information
type Damage struct {
	TargetId     int32
	Effects      []byte
	DamageAmount uint16
	Flags        byte
	Killed       bool
	ArmorPierce  bool
	Laser        bool
	BulletId     uint16
	ObjectId     int32
}

// Constants for damage flags
const (
	DamageFlagKill        = 1
	DamageFlagArmorPierce = 2
	DamageFlagLaser       = 4
)

// Type returns the packet type for Damage
func (p *Damage) Type() interfaces.PacketType {
	return interfaces.Damage
}

// Read reads the packet data from the provided reader
func (p *Damage) Read(r interfaces.Reader) error {
	var err error

	// Read TargetId
	p.TargetId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Effects
	effectsLength, err := r.ReadByte()
	if err != nil {
		return err
	}

	p.Effects = make([]byte, effectsLength)
	for i := 0; i < int(effectsLength); i++ {
		p.Effects[i], err = r.ReadByte()
		if err != nil {
			return err
		}
	}

	// Read DamageAmount
	p.DamageAmount, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	// Read Flags
	p.Flags, err = r.ReadByte()
	if err != nil {
		return err
	}

	// Set flag-based booleans
	p.Killed = (p.Flags & DamageFlagKill) > 0
	p.ArmorPierce = (p.Flags & DamageFlagArmorPierce) > 0
	p.Laser = (p.Flags & DamageFlagLaser) > 0

	// Read BulletId
	p.BulletId, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	// Read ObjectId
	p.ObjectId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *Damage) Write(w interfaces.Writer) error {
	var err error

	// Write TargetId
	err = w.WriteInt32(p.TargetId)
	if err != nil {
		return err
	}

	// Write Effects
	err = w.WriteByte(byte(len(p.Effects)))
	if err != nil {
		return err
	}

	for _, effect := range p.Effects {
		err = w.WriteByte(effect)
		if err != nil {
			return err
		}
	}

	// Write DamageAmount
	err = w.WriteUInt16(p.DamageAmount)
	if err != nil {
		return err
	}

	// Write Flags
	err = w.WriteByte(p.Flags)
	if err != nil {
		return err
	}

	// Write BulletId
	err = w.WriteUInt16(p.BulletId)
	if err != nil {
		return err
	}

	// Write ObjectId
	err = w.WriteInt32(p.ObjectId)
	if err != nil {
		return err
	}

	return nil
}
