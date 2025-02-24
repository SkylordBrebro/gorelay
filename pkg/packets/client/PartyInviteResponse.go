package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// AcceptDecline represents the response to a party invite
type AcceptDecline byte

const (
	// Decline represents declining the party invite
	Decline AcceptDecline = iota
	// Accept represents accepting the party invite
	Accept
)

// PartyInviteResponse represents a packet for responding to party invites
type PartyInviteResponse struct {
	*packets.BasePacket
	PartyID      uint32
	AcceptInvite AcceptDecline
}

// NewPartyInviteResponse creates a new PartyInviteResponse packet
func NewPartyInviteResponse() *PartyInviteResponse {
	return &PartyInviteResponse{
		BasePacket: packets.NewPacket(interfaces.PartyInviteResponse, 253), // 253 is the unsigned byte equivalent of -3
	}
}

// Type returns the packet type
func (p *PartyInviteResponse) Type() interfaces.PacketType {
	return interfaces.PartyInviteResponse
}

// Read reads the packet data from a PacketReader
func (p *PartyInviteResponse) Read(r *packets.PacketReader) error {
	var err error
	p.PartyID, err = r.ReadUInt32()
	if err != nil {
		return err
	}
	acceptByte, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.AcceptInvite = AcceptDecline(acceptByte)
	return nil
}

// Write writes the packet data to a PacketWriter
func (p *PartyInviteResponse) Write(w *packets.PacketWriter) error {
	if err := w.WriteUInt32(p.PartyID); err != nil {
		return err
	}
	return w.WriteByte(byte(p.AcceptInvite))
}
