package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// CancelTrade represents a packet for canceling a trade
type CancelTrade struct {
	*packets.BasePacket
}

// NewCancelTrade creates a new CancelTrade packet
func NewCancelTrade() *CancelTrade {
	return &CancelTrade{
		BasePacket: packets.NewPacket(interfaces.CancelTrade, byte(interfaces.CancelTrade)),
	}
}

// Type returns the packet type
func (p *CancelTrade) Type() interfaces.PacketType {
	return interfaces.CancelTrade
}

// Read reads the packet data from a PacketReader
func (p *CancelTrade) Read(r *packets.PacketReader) error {
	return nil
}

// Write writes the packet data to a PacketWriter
func (p *CancelTrade) Write(w *packets.PacketWriter) error {
	return nil
}
