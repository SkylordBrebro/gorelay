package server

import (
	"gorelay/pkg/packets/interfaces"
)

// EnemyShoot represents the server packet for enemy shooting information
type EnemyShoot struct {
	BulletId   uint16
	OwnerId    int32
	BulletType byte
	Location   Location
	Angle      float32
	Damage     int16
	NumShots   byte
	AngleInc   float32
}

// Type returns the packet type for EnemyShoot
func (p *EnemyShoot) Type() interfaces.PacketType {
	return interfaces.EnemyShoot
}

// Read reads the packet data from the provided reader
func (p *EnemyShoot) Read(r interfaces.Reader) error {
	var err error

	// Read BulletId
	p.BulletId, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	// Read OwnerId
	p.OwnerId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read BulletType
	p.BulletType, err = r.ReadByte()
	if err != nil {
		return err
	}

	// Read Location
	p.Location, err = NewLocation(r)
	if err != nil {
		return err
	}

	// Read Angle
	p.Angle, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	// Read Damage
	p.Damage, err = r.ReadInt16()
	if err != nil {
		return err
	}

	// Check if we have more data to read
	if r.RemainingBytes() > 0 {
		// Read NumShots
		p.NumShots, err = r.ReadByte()
		if err != nil {
			return err
		}

		// Read AngleInc
		p.AngleInc, err = r.ReadFloat32()
		if err != nil {
			return err
		}
	} else {
		// Default values if not present
		p.NumShots = 1
		p.AngleInc = 0.0
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *EnemyShoot) Write(w interfaces.Writer) error {
	var err error

	// Write BulletId
	err = w.WriteUInt16(p.BulletId)
	if err != nil {
		return err
	}

	// Write OwnerId
	err = w.WriteInt32(p.OwnerId)
	if err != nil {
		return err
	}

	// Write BulletType
	err = w.WriteByte(p.BulletType)
	if err != nil {
		return err
	}

	// Write Location
	err = p.Location.Write(w)
	if err != nil {
		return err
	}

	// Write Angle
	err = w.WriteFloat32(p.Angle)
	if err != nil {
		return err
	}

	// Write Damage
	err = w.WriteInt16(p.Damage)
	if err != nil {
		return err
	}

	// Only write NumShots and AngleInc if NumShots is not 1
	if p.NumShots != 1 {
		// Write NumShots
		err = w.WriteByte(p.NumShots)
		if err != nil {
			return err
		}

		// Write AngleInc
		err = w.WriteFloat32(p.AngleInc)
		if err != nil {
			return err
		}
	}

	return nil
}
