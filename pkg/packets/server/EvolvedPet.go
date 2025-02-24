package server

import (
	"gorelay/pkg/packets/interfaces"
)

// EvolvedPet represents the server packet for pet evolution information
type EvolvedPet struct {
	PetId       int32
	InitialSkin int32
	FinalSkin   int32
}

// Type returns the packet type for EvolvedPet
func (p *EvolvedPet) Type() interfaces.PacketType {
	return interfaces.EvolvedPet
}

// Read reads the packet data from the provided reader
func (p *EvolvedPet) Read(r interfaces.Reader) error {
	var err error

	p.PetId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.InitialSkin, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.FinalSkin, err = r.ReadInt32()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *EvolvedPet) Write(w interfaces.Writer) error {
	var err error

	err = w.WriteInt32(p.PetId)
	if err != nil {
		return err
	}

	err = w.WriteInt32(p.InitialSkin)
	if err != nil {
		return err
	}

	err = w.WriteInt32(p.FinalSkin)
	if err != nil {
		return err
	}

	return nil
}
