package server

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// ServerPlayerShoot represents an incoming server player shoot packet
type ServerPlayerShoot struct {
	BulletId      int32
	OwnerId       int32
	ContainerType int32
	StartingPos   *dataobjects.Location
	Angle         float32
	Damage        int16
}

// Type returns the packet type for ServerPlayerShoot
func (p *ServerPlayerShoot) Type() interfaces.PacketType {
	return interfaces.ServerPlayerShoot
}

// Read reads the packet data from the provided reader
func (p *ServerPlayerShoot) Read(r interfaces.Reader) error {
	var err error

	p.BulletId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.OwnerId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.ContainerType, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.StartingPos = &dataobjects.Location{}
	err = p.StartingPos.Read(r)
	if err != nil {
		return err
	}

	p.Angle, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Damage, err = r.ReadInt16()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *ServerPlayerShoot) Write(w interfaces.Writer) error {
	var err error

	err = w.WriteInt32(p.BulletId)
	if err != nil {
		return err
	}

	err = w.WriteInt32(p.OwnerId)
	if err != nil {
		return err
	}

	err = w.WriteInt32(p.ContainerType)
	if err != nil {
		return err
	}

	err = p.StartingPos.Write(w)
	if err != nil {
		return err
	}

	err = w.WriteFloat32(p.Angle)
	if err != nil {
		return err
	}

	err = w.WriteInt16(p.Damage)
	if err != nil {
		return err
	}

	return nil
}

func (p *ServerPlayerShoot) ID() int32 {
	return int32(interfaces.ServerPlayerShoot)
}
