package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// PartyJoinRequest represents a packet for party join requests
type PartyJoinRequest struct {
	*packets.BasePacket
	PlayerID uint32
	Unknown1 byte
}

// NewPartyJoinRequest creates a new PartyJoinRequest packet
func NewPartyJoinRequest() *PartyJoinRequest {
	return &PartyJoinRequest{
		BasePacket: packets.NewPacket(interfaces.PartyJoinRequest, 252), // 252 is the unsigned byte equivalent of -4
	}
}

// Type returns the packet type
func (p *PartyJoinRequest) Type() interfaces.PacketType {
	return interfaces.PartyJoinRequest
}

// Read reads the packet data from a PacketReader
func (p *PartyJoinRequest) Read(r *packets.PacketReader) error {
	var err error
	p.PlayerID, err = r.ReadUInt32()
	if err != nil {
		return err
	}
	p.Unknown1, err = r.ReadByte()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *PartyJoinRequest) Write(w *packets.PacketWriter) error {
	if err := w.WriteUInt32(p.PlayerID); err != nil {
		return err
	}
	return w.WriteByte(p.Unknown1)
}
