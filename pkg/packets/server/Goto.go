package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Goto represents the server packet for teleportation information
type Goto struct {
	ObjectId int32
	Location Location
	Unknown  int32
}

// Type returns the packet type for Goto
func (p *Goto) Type() interfaces.PacketType {
	return interfaces.Goto
}

// Read reads the packet data from the provided reader
func (p *Goto) Read(r interfaces.Reader) error {
	var err error

	// Read ObjectId
	p.ObjectId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Location
	p.Location, err = NewLocation(r)
	if err != nil {
		return err
	}

	// Read Unknown
	p.Unknown, err = r.ReadInt32()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *Goto) Write(w interfaces.Writer) error {
	var err error

	// Write ObjectId
	err = w.WriteInt32(p.ObjectId)
	if err != nil {
		return err
	}

	// Write Location
	err = p.Location.Write(w)
	if err != nil {
		return err
	}

	// Write Unknown
	err = w.WriteInt32(p.Unknown)
	if err != nil {
		return err
	}

	return nil
}
