package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// InventorySwap represents a packet for swapping inventory items
type InventorySwap struct {
	*packets.BasePacket
	Time        int32
	Position    *Location
	SlotObject1 *dataobjects.SlotObject
	SlotObject2 *dataobjects.SlotObject
}

// NewInventorySwap creates a new InventorySwap packet
func NewInventorySwap() *InventorySwap {
	return &InventorySwap{
		BasePacket:  packets.NewPacket(interfaces.InventorySwap, byte(interfaces.InventorySwap)),
		SlotObject1: dataobjects.NewSlotObject(),
		SlotObject2: dataobjects.NewSlotObject(),
	}
}

// Type returns the packet type
func (p *InventorySwap) Type() interfaces.PacketType {
	return interfaces.InventorySwap
}

// Read reads the packet data from a PacketReader
func (p *InventorySwap) Read(r *packets.PacketReader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Position, err = NewLocation(r)
	if err != nil {
		return err
	}
	if err = p.SlotObject1.Read(r); err != nil {
		return err
	}
	return p.SlotObject2.Read(r)
}

// Write writes the packet data to a PacketWriter
func (p *InventorySwap) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	if err := p.Position.Write(w); err != nil {
		return err
	}
	if err := p.SlotObject1.Write(w); err != nil {
		return err
	}
	return p.SlotObject2.Write(w)
}
