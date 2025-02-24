package server

import (
	"gorelay/pkg/packets/interfaces"
)

// SlotObject represents an item slot in the game
type SlotObject struct {
	ObjectId   int32
	SlotId     int
	ObjectType int32
}

// NewSlotObject creates a new SlotObject by reading from a packet reader
func NewSlotObject(r interfaces.Reader) (SlotObject, error) {
	var slot SlotObject
	var err error

	slot.ObjectId, err = r.ReadInt32()
	if err != nil {
		return slot, err
	}

	slot.SlotId, err = r.ReadCompressedInt()
	if err != nil {
		return slot, err
	}

	slot.ObjectType, err = r.ReadInt32()
	return slot, err
}

// Write writes the slot object data to a packet writer
func (s *SlotObject) Write(w interfaces.Writer) error {
	var err error

	err = w.WriteInt32(s.ObjectId)
	if err != nil {
		return err
	}

	err = w.WriteCompressedInt(s.SlotId)
	if err != nil {
		return err
	}

	err = w.WriteInt32(s.ObjectType)
	return err
}

// ForgeResult represents the server packet for forge operation results
type ForgeResult struct {
	Success         bool
	DismantledSlots []SlotObject
}

// Type returns the packet type for ForgeResult
func (p *ForgeResult) Type() interfaces.PacketType {
	return interfaces.ForgeResult
}

// Read reads the packet data from the provided reader
func (p *ForgeResult) Read(r interfaces.Reader) error {
	var err error

	// Read Success
	p.Success, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read DismantledSlots
	slotsLength, err := r.ReadByte()
	if err != nil {
		return err
	}

	p.DismantledSlots = make([]SlotObject, slotsLength)
	for i := 0; i < int(slotsLength); i++ {
		p.DismantledSlots[i], err = NewSlotObject(r)
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *ForgeResult) Write(w interfaces.Writer) error {
	var err error

	// Write Success
	err = w.WriteBool(p.Success)
	if err != nil {
		return err
	}

	// Write DismantledSlots length
	err = w.WriteByte(byte(len(p.DismantledSlots)))
	if err != nil {
		return err
	}

	// Write DismantledSlots
	for _, slot := range p.DismantledSlots {
		err = slot.Write(w)
		if err != nil {
			return err
		}
	}

	return nil
}
