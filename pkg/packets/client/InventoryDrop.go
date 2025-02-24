package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// InventoryDrop represents a packet for dropping inventory items
type InventoryDrop struct {
	*packets.BasePacket
	Slot    *dataobjects.SlotObject
	Unknown bool
}

// NewInventoryDrop creates a new InventoryDrop packet
func NewInventoryDrop() *InventoryDrop {
	return &InventoryDrop{
		BasePacket: packets.NewPacket(interfaces.InventoryDrop, byte(interfaces.InventoryDrop)),
		Slot:       dataobjects.NewSlotObject(),
	}
}

// Type returns the packet type
func (p *InventoryDrop) Type() interfaces.PacketType {
	return interfaces.InventoryDrop
}

// Read reads the packet data from a PacketReader
func (p *InventoryDrop) Read(r *packets.PacketReader) error {
	var err error
	if err = p.Slot.Read(r); err != nil {
		return err
	}
	p.Unknown, err = r.ReadBool()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *InventoryDrop) Write(w *packets.PacketWriter) error {
	if err := p.Slot.Write(w); err != nil {
		return err
	}
	return w.WriteBool(p.Unknown)
}
