package server

import (
	"gorelay/pkg/packets/interfaces"
)

// DeletePet represents the server packet for pet deletion
type DeletePet struct {
	PetId int32
}

// Type returns the packet type for DeletePet
func (p *DeletePet) Type() interfaces.PacketType {
	return interfaces.DeletePet
}

// Read reads the packet data from the provided reader
func (p *DeletePet) Read(r interfaces.Reader) error {
	var err error
	p.PetId, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the provided writer
func (p *DeletePet) Write(w interfaces.Writer) error {
	return w.WriteInt32(p.PetId)
}
