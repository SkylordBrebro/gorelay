package server

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// AllyShoot represents an incoming ally shoot packet from the server
type AllyShoot struct {
	BulletId      int32
	OwnerId       int32
	ContainerType int32
	StartingPos   *dataobjects.Location
	Angle         float32
	Damage        int16
}

// Type returns the packet type for AllyShoot
func (p *AllyShoot) Type() interfaces.PacketType {
	return interfaces.AllyShoot
}

// Read reads the packet data from the provided reader
func (p *AllyShoot) Read(r interfaces.Reader) error {
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
func (p *AllyShoot) Write(w interfaces.Writer) error {
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

func (p *AllyShoot) ID() int32 {
	return int32(interfaces.AllyShoot)
}

// String returns a string representation of the packet
func (p *AllyShoot) String() string {
	return "AllyShoot"
}

// HasNulls checks if any fields in the packet are null
func (p *AllyShoot) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *AllyShoot) Structure() string {
	return "AllyShoot"
}
