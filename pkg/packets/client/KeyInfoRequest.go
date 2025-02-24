package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// KeyInfoRequest represents a packet for requesting key information
type KeyInfoRequest struct {
	*packets.BasePacket
	ItemID int32
}

// NewKeyInfoRequest creates a new KeyInfoRequest packet
func NewKeyInfoRequest() *KeyInfoRequest {
	return &KeyInfoRequest{
		BasePacket: packets.NewPacket(interfaces.KeyInfoRequest, byte(interfaces.KeyInfoRequest)),
	}
}

// Type returns the packet type
func (p *KeyInfoRequest) Type() interfaces.PacketType {
	return interfaces.KeyInfoRequest
}

// Read reads the packet data from a PacketReader
func (p *KeyInfoRequest) Read(r *packets.PacketReader) error {
	var err error
	p.ItemID, err = r.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *KeyInfoRequest) Write(w *packets.PacketWriter) error {
	return w.WriteInt32(p.ItemID)
}
