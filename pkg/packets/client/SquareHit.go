package client

import (
	"gorelay/pkg/packets/interfaces"
)

// SquareHit represents a client-side square hit packet
type SquareHit struct {
	Time     int32
	BulletId uint16
	ObjectId int32
}

// Type returns the packet type for SquareHit
func (p *SquareHit) Type() interfaces.PacketType {
	return interfaces.SquareHit
}

// Read reads the packet data from the given reader
func (p *SquareHit) Read(r interfaces.Reader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.BulletId, err = r.ReadUInt16()
	if err != nil {
		return err
	}
	p.ObjectId, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the given writer
func (p *SquareHit) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	if err := w.WriteUInt16(p.BulletId); err != nil {
		return err
	}
	return w.WriteInt32(p.ObjectId)
}
