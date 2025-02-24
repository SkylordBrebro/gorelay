package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// BuyEmote represents an emote purchase packet
type BuyEmote struct {
	*packets.BasePacket
	EmoteID int32
}

// NewBuyEmote creates a new BuyEmote packet
func NewBuyEmote() *BuyEmote {
	return &BuyEmote{
		BasePacket: packets.NewPacket(interfaces.BuyEmote, byte(interfaces.BuyEmote)),
	}
}

// Type returns the packet type
func (p *BuyEmote) Type() interfaces.PacketType {
	return interfaces.BuyEmote
}

// Read reads the packet data from the reader
func (p *BuyEmote) Read(r *packets.PacketReader) error {
	var err error
	p.EmoteID, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the writer
func (p *BuyEmote) Write(w *packets.PacketWriter) error {
	return w.WriteInt32(p.EmoteID)
}
