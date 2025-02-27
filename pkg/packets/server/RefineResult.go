package server

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// RefineResult represents the server packet for refine results
type RefineResult struct {
	Success bool
	Slot    *dataobjects.SlotObject
	KeyMods string
}

// Type returns the packet type for RefineResult
func (p *RefineResult) Type() interfaces.PacketType {
	return interfaces.RefineResult
}

// Read reads the packet data from the provided reader
func (p *RefineResult) Read(r interfaces.Reader) error {
	var err error

	// Read Success
	p.Success, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read Slot
	p.Slot = dataobjects.NewSlotObject()
	err = p.Slot.Read(r)
	if err != nil {
		return err
	}

	// Read KeyMods
	p.KeyMods, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *RefineResult) Write(w interfaces.Writer) error {
	var err error

	// Write Success
	err = w.WriteBool(p.Success)
	if err != nil {
		return err
	}

	// Write Slot
	err = p.Slot.Write(w)
	if err != nil {
		return err
	}

	// Write KeyMods
	err = w.WriteString(p.KeyMods)
	if err != nil {
		return err
	}

	return nil
}

func (p *RefineResult) ID() int32 {
	return int32(interfaces.RefineResult)
}