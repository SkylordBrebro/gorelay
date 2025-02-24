package server

import (
	"gorelay/pkg/packets/interfaces"
)

// HatchPet represents the server packet for pet hatching information
type HatchPet struct {
	PetName   string
	PetSkinId int32
	ItemType  int32
}

// Type returns the packet type for HatchPet
func (p *HatchPet) Type() interfaces.PacketType {
	return interfaces.HatchPet
}

// Read reads the packet data from the provided reader
func (p *HatchPet) Read(r interfaces.Reader) error {
	var err error

	// Read PetName
	p.PetName, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read PetSkinId
	p.PetSkinId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read ItemType
	p.ItemType, err = r.ReadInt32()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *HatchPet) Write(w interfaces.Writer) error {
	var err error

	// Write PetName
	err = w.WriteString(p.PetName)
	if err != nil {
		return err
	}

	// Write PetSkinId
	err = w.WriteInt32(p.PetSkinId)
	if err != nil {
		return err
	}

	// Write ItemType
	err = w.WriteInt32(p.ItemType)
	if err != nil {
		return err
	}

	return nil
}
