package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Location represents a position in the game world
type Location struct {
	X float32
	Y float32
}

// NewLocation creates a new Location by reading from a packet reader
func NewLocation(r interfaces.Reader) (Location, error) {
	var loc Location
	var err error

	loc.X, err = r.ReadFloat32()
	if err != nil {
		return loc, err
	}

	loc.Y, err = r.ReadFloat32()
	return loc, err
}

// Write writes the location data to a packet writer
func (l *Location) Write(w interfaces.Writer) error {
	err := w.WriteFloat32(l.X)
	if err != nil {
		return err
	}

	return w.WriteFloat32(l.Y)
}

// DrawDebugArrow represents the server packet for debug arrow drawing
type DrawDebugArrow struct {
	ID       uint32
	StartLoc Location
	EndLoc   Location
	Lifetime float32
	Color    float32
}

// Type returns the packet type for DrawDebugArrow
func (p *DrawDebugArrow) Type() interfaces.PacketType {
	return interfaces.DrawDebugArrow
}

// Read reads the packet data from the provided reader
func (p *DrawDebugArrow) Read(r interfaces.Reader) error {
	var err error

	// Read ID
	p.ID, err = r.ReadUInt32()
	if err != nil {
		return err
	}

	// Read StartLoc
	p.StartLoc, err = NewLocation(r)
	if err != nil {
		return err
	}

	// Read EndLoc
	p.EndLoc, err = NewLocation(r)
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
func (p *DrawDebugArrow) Write(w interfaces.Writer) error {
	var err error

	// Write ID
	err = w.WriteUInt32(p.ID)
	if err != nil {
		return err
	}

	// Write StartLoc
	err = p.StartLoc.Write(w)
	if err != nil {
		return err
	}

	// Write EndLoc
	err = p.EndLoc.Write(w)
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
