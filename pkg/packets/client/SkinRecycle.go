package client

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// SkinRecycle represents a client-side skin recycling packet
type SkinRecycle struct {
	Item *dataobjects.SlotObject
}

// Type returns the packet type for SkinRecycle
func (p *SkinRecycle) Type() interfaces.PacketType {
	return interfaces.SkinRecycle
}

// Read reads the packet data from the given reader
func (p *SkinRecycle) Read(r interfaces.Reader) error {
	p.Item = dataobjects.NewSlotObject()
	return p.Item.Read(r)
}

// Write writes the packet data to the given writer
func (p *SkinRecycle) Write(w interfaces.Writer) error {
	return p.Item.Write(w)
}
