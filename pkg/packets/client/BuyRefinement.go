package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// BuyRefinement represents a packet for buying refinements
type BuyRefinement struct {
	*packets.BasePacket
	Slot   *dataobjects.SlotObject
	Action int16
}

// NewBuyRefinement creates a new BuyRefinement packet
func NewBuyRefinement() *BuyRefinement {
	return &BuyRefinement{
		BasePacket: packets.NewPacket(interfaces.BuyRefinement, byte(interfaces.BuyRefinement)),
		Slot:       dataobjects.NewSlotObject(),
	}
}

// Type returns the packet type
func (p *BuyRefinement) Type() interfaces.PacketType {
	return interfaces.BuyRefinement
}

// Read reads the packet data from a PacketReader
func (p *BuyRefinement) Read(r *packets.PacketReader) error {
	var err error
	if err = p.Slot.Read(r); err != nil {
		return err
	}
	p.Action, err = r.ReadInt16()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *BuyRefinement) Write(w *packets.PacketWriter) error {
	if err := p.Slot.Write(w); err != nil {
		return err
	}
	return w.WriteInt16(p.Action)
}
