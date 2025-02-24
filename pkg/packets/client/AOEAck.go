package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// Location represents a position in the game world
type Location struct {
	X float32
	Y float32
}

// NewLocation creates a new Location from a packet reader
func NewLocation(r *packets.PacketReader) (*Location, error) {
	x, err := r.ReadFloat32()
	if err != nil {
		return nil, err
	}

	y, err := r.ReadFloat32()
	if err != nil {
		return nil, err
	}

	return &Location{X: x, Y: y}, nil
}

// Write writes the location to a packet writer
func (l *Location) Write(w *packets.PacketWriter) error {
	if err := w.WriteFloat32(l.X); err != nil {
		return err
	}
	return w.WriteFloat32(l.Y)
}

// AOEAck represents an area of effect acknowledgment packet
type AOEAck struct {
	*packets.BasePacket
	Time     int32
	Position *Location
}

// NewAOEAck creates a new AOEAck packet
func NewAOEAck() *AOEAck {
	return &AOEAck{
		BasePacket: packets.NewPacket(interfaces.AOEAck, byte(interfaces.AOEAck)),
	}
}

// Type returns the packet type
func (p *AOEAck) Type() interfaces.PacketType {
	return interfaces.AOEAck
}

// Read reads the packet data from the reader
func (p *AOEAck) Read(r *packets.PacketReader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Position, err = NewLocation(r)
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the writer
func (p *AOEAck) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}

	return p.Position.Write(w)
}
