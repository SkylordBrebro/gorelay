package server

import (
	"gorelay/pkg/packets/interfaces"
)

// DrawDebugShape represents the server packet for debug shape drawing
type DrawDebugShape struct {
	ID        uint32
	Location  Location
	ShapeType byte
	Lifetime  float32
	Color     float32
}

// Type returns the packet type for DrawDebugShape
func (p *DrawDebugShape) Type() interfaces.PacketType {
	return interfaces.DrawDebugShape
}

// Read reads the packet data from the provided reader
func (p *DrawDebugShape) Read(r interfaces.Reader) error {
	var err error

	// Read ID
	p.ID, err = r.ReadUInt32()
	if err != nil {
		return err
	}

	// Read Location
	p.Location, err = NewLocation(r)
	if err != nil {
		return err
	}

	// Read ShapeType
	p.ShapeType, err = r.ReadByte()
	if err != nil {
		return err
	}

	// Read Lifetime
	p.Lifetime, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	// Read Color
	p.Color, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *DrawDebugShape) Write(w interfaces.Writer) error {
	var err error

	// Write ID
	err = w.WriteUInt32(p.ID)
	if err != nil {
		return err
	}

	// Write Location
	err = p.Location.Write(w)
	if err != nil {
		return err
	}

	// Write ShapeType
	err = w.WriteByte(p.ShapeType)
	if err != nil {
		return err
	}

	// Write Lifetime
	err = w.WriteFloat32(p.Lifetime)
	if err != nil {
		return err
	}

	// Write Color
	err = w.WriteFloat32(p.Color)
	if err != nil {
		return err
	}

	return nil
}
