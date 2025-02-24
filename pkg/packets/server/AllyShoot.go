package server

import (
	"gorelay/pkg/packets/interfaces"
)

// AllyShoot represents a server-side ally shoot packet
type AllyShoot struct {
	BulletId      uint16
	OwnerId       int32
	ContainerType int16
	ProjPosition  int8
	Angle         float32
	UnknownSbyte0 int8
	AttackType    int8
	UnknownInt0   int32
}

// Type returns the packet type for AllyShoot
func (p *AllyShoot) Type() interfaces.PacketType {
	return interfaces.AllyShoot
}

// ID returns the packet ID
func (p *AllyShoot) ID() int32 {
	return int32(interfaces.AllyShoot)
}

// Read reads the packet data from the given reader
func (p *AllyShoot) Read(r interfaces.Reader) error {
	var err error
	p.BulletId, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	p.OwnerId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.ContainerType, err = r.ReadInt16()
	if err != nil {
		return err
	}

	projPos, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.ProjPosition = int8(projPos)

	p.Angle, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	unknownByte0, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.UnknownSbyte0 = int8(unknownByte0)

	attackType, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.AttackType = int8(attackType)

	// Check if there are 4 more bytes to read
	if r.RemainingBytes() >= 4 {
		p.UnknownInt0, err = r.ReadInt32()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the given writer
func (p *AllyShoot) Write(w interfaces.Writer) error {
	if err := w.WriteUInt16(p.BulletId); err != nil {
		return err
	}

	if err := w.WriteInt32(p.OwnerId); err != nil {
		return err
	}

	if err := w.WriteInt16(p.ContainerType); err != nil {
		return err
	}

	if err := w.WriteByte(byte(p.ProjPosition)); err != nil {
		return err
	}

	if err := w.WriteFloat32(p.Angle); err != nil {
		return err
	}

	if err := w.WriteByte(byte(p.UnknownSbyte0)); err != nil {
		return err
	}

	if err := w.WriteByte(byte(p.AttackType)); err != nil {
		return err
	}

	// Only write the UnknownInt0 if it's not 1
	if p.UnknownInt0 != 1 {
		if err := w.WriteInt32(p.UnknownInt0); err != nil {
			return err
		}
	}

	return nil
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
