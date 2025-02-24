package server

import (
	"gorelay/pkg/packets/interfaces"
)

// PetYardUpdate represents the server packet for pet yard updates
type PetYardUpdate struct {
	TypeId int32
}

// Type returns the packet type for PetYardUpdate
func (p *PetYardUpdate) Type() interfaces.PacketType {
	return interfaces.PetYardUpdate
}

// Read reads the packet data from the provided reader
func (p *PetYardUpdate) Read(r interfaces.Reader) error {
	var err error
	p.TypeId, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the provided writer
func (p *PetYardUpdate) Write(w interfaces.Writer) error {
	return w.WriteInt32(p.TypeId)
}
