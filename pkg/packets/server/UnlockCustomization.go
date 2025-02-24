package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Customization type constants
const (
	CustomizationTypeSkin       = 1
	CustomizationTypePetSkin    = 2
	CustomizationTypeTitle      = 3
	CustomizationTypeGravestone = 4
	CustomizationTypeEmote      = 5
)

// UnlockCustomization represents the server packet for unlocking customizations
type UnlockCustomization struct {
	UnlockType int8
	SkinType   int8
	ItemType   int32
	CostType   int32
}

// Type returns the packet type for UnlockCustomization
func (p *UnlockCustomization) Type() interfaces.PacketType {
	return interfaces.UnlockCustomization
}

// Read reads the packet data from the provided reader
func (p *UnlockCustomization) Read(r interfaces.Reader) error {
	var err error

	// Read UnlockType
	unlockTypeByte, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.UnlockType = int8(unlockTypeByte)

	// Read SkinType
	skinTypeByte, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.SkinType = int8(skinTypeByte)

	// Read ItemType
	p.ItemType, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read CostType
	p.CostType, err = r.ReadInt32()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *UnlockCustomization) Write(w interfaces.Writer) error {
	var err error

	// Write UnlockType
	err = w.WriteByte(byte(p.UnlockType))
	if err != nil {
		return err
	}

	// Write SkinType
	err = w.WriteByte(byte(p.SkinType))
	if err != nil {
		return err
	}

	// Write ItemType
	err = w.WriteInt32(p.ItemType)
	if err != nil {
		return err
	}

	// Write CostType
	err = w.WriteInt32(p.CostType)
	if err != nil {
		return err
	}

	return nil
}
