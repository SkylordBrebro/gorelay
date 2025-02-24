package server

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// SkinRecycleResponse represents the server packet for skin recycle responses
type SkinRecycleResponse struct {
	Success bool
	Item    *dataobjects.SlotObject
}

// Type returns the packet type for SkinRecycleResponse
func (p *SkinRecycleResponse) Type() interfaces.PacketType {
	return interfaces.SkinRecycleResponse
}

// Read reads the packet data from the provided reader
func (p *SkinRecycleResponse) Read(r interfaces.Reader) error {
	var err error

	// Read Success
	p.Success, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read Item
	p.Item = dataobjects.NewSlotObject()
	err = p.Item.Read(r)
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *SkinRecycleResponse) Write(w interfaces.Writer) error {
	var err error

	// Write Success
	err = w.WriteBool(p.Success)
	if err != nil {
		return err
	}

	// Write Item
	err = p.Item.Write(w)
	if err != nil {
		return err
	}

	return nil
}
