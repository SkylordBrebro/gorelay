package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// BuyItem represents an item purchase packet
type BuyItem struct {
	*packets.BasePacket
	ItemIDs []int32
}

// NewBuyItem creates a new BuyItem packet
func NewBuyItem() *BuyItem {
	return &BuyItem{
		BasePacket: packets.NewPacket(interfaces.BuyItem, byte(interfaces.BuyItem)),
		ItemIDs:    make([]int32, 0),
	}
}

// Type returns the packet type
func (p *BuyItem) Type() interfaces.PacketType {
	return interfaces.BuyItem
}

// Read reads the packet data from the reader
func (p *BuyItem) Read(r *packets.PacketReader) error {
	// Read array length
	length, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read item IDs
	p.ItemIDs = make([]int32, length)
	for i := 0; i < int(length); i++ {
		p.ItemIDs[i], err = r.ReadInt32()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the writer
func (p *BuyItem) Write(w *packets.PacketWriter) error {
	// Write array length
	if err := w.WriteInt16(int16(len(p.ItemIDs))); err != nil {
		return err
	}

	// Write item IDs
	for _, id := range p.ItemIDs {
		if err := w.WriteInt32(id); err != nil {
			return err
		}
	}

	return nil
}
