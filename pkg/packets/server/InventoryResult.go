package server

import (
	"gorelay/pkg/packets/interfaces"
)

// InventoryResult represents the server packet for inventory operation results
type InventoryResult struct {
	Result      bool
	Unknown     byte
	SlotObject1 SlotObject
	SlotObject2 SlotObject
}

// Type returns the packet type for InventoryResult
func (p *InventoryResult) Type() interfaces.PacketType {
	return interfaces.InventoryResult
}

// Read reads the packet data from the provided reader
func (p *InventoryResult) Read(r interfaces.Reader) error {
	var err error

	// Read Result
	p.Result, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read Unknown
	p.Unknown, err = r.ReadByte()
	if err != nil {
		return err
	}

	// Read SlotObject1
	p.SlotObject1, err = NewSlotObject(r)
	if err != nil {
		return err
	}

	// Read SlotObject2
	p.SlotObject2, err = NewSlotObject(r)
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *InventoryResult) Write(w interfaces.Writer) error {
	var err error

	// Write Result
	err = w.WriteBool(p.Result)
	if err != nil {
		return err
	}

	// Write Unknown
	err = w.WriteByte(p.Unknown)
	if err != nil {
		return err
	}

	// Write SlotObject1
	err = p.SlotObject1.Write(w)
	if err != nil {
		return err
	}

	// Write SlotObject2
	err = p.SlotObject2.Write(w)
	if err != nil {
		return err
	}

	return nil
}

func (p *InventoryResult) ID() int32 {
	return int32(interfaces.InventoryResult)
}