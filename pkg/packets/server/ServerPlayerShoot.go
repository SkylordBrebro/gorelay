package server

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// ServerPlayerShoot represents the server packet for player shooting
type ServerPlayerShoot struct {
	BulletId        int16
	OwnerId         int32
	ContainerType   int32
	StartingLoc     *dataobjects.Location
	Angle           float32
	Damage          int16
	TargetId        int32
	ProjectileCount byte
	UnknownByte1    byte
	AngleIncrement  float32
}

// Type returns the packet type for ServerPlayerShoot
func (p *ServerPlayerShoot) Type() interfaces.PacketType {
	return interfaces.ServerPlayerShoot
}

// Read reads the packet data from the provided reader
func (p *ServerPlayerShoot) Read(r interfaces.Reader) error {
	var err error

	// Read BulletId
	p.BulletId, err = r.ReadInt16()
	if err != nil {
		return err
	}

	// Read OwnerId
	p.OwnerId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read ContainerType
	p.ContainerType, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read StartingLoc
	p.StartingLoc = dataobjects.NewLocation()
	err = p.StartingLoc.Read(r)
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

	// Read TargetId
	p.TargetId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Initialize default values
	p.UnknownByte1 = 1
	p.AngleIncrement = 0

	// Check if there are more bytes to read for optional fields
	if r.RemainingBytes() > 0 {
		// Try to read ProjectileCount
		p.ProjectileCount, err = r.ReadByte()
		if err != nil {
			return err
		}

		// Check if there are more bytes to read for UnknownByte1
		if r.RemainingBytes() > 0 {
			p.UnknownByte1, err = r.ReadByte()
			if err != nil {
				return err
			}

			// Check if there are more bytes to read for AngleIncrement
			if r.RemainingBytes() >= 4 {
				p.AngleIncrement, err = r.ReadFloat32()
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *ServerPlayerShoot) Write(w interfaces.Writer) error {
	var err error

	// Write BulletId
	err = w.WriteInt16(p.BulletId)
	if err != nil {
		return err
	}

	// Write OwnerId
	err = w.WriteInt32(p.OwnerId)
	if err != nil {
		return err
	}

	// Write ContainerType
	err = w.WriteInt32(p.ContainerType)
	if err != nil {
		return err
	}

	// Write StartingLoc
	err = p.StartingLoc.Write(w)
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

	// Write TargetId
	err = w.WriteInt32(p.TargetId)
	if err != nil {
		return err
	}

	// Write optional fields if needed
	if p.ProjectileCount > 0 {
		err = w.WriteByte(p.ProjectileCount)
		if err != nil {
			return err
		}
	}

	if p.UnknownByte1 != 1 {
		err = w.WriteByte(p.UnknownByte1)
		if err != nil {
			return err
		}
	}

	if p.AngleIncrement != 0 {
		err = w.WriteFloat32(p.AngleIncrement)
		if err != nil {
			return err
		}
	}

	return nil
}
