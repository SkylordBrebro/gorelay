package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// PartyActionID represents different party action types
type PartyActionID byte

const (
	// PartyActionNone represents no action
	PartyActionNone PartyActionID = iota
	// PartyActionFailed represents failed action
	PartyActionFailed
	// PartyActionKicked represents kicked from party
	PartyActionKicked
	// PartyActionKickNotFound represents kick target not found
	PartyActionKickNotFound
	// PartyActionPromotedToLeader represents promoted to leader
	PartyActionPromotedToLeader
	// PartyActionPromoteNotFound represents promote target not found
	PartyActionPromoteNotFound
	// PartyActionLeftParty represents left party
	PartyActionLeftParty
	// PartyActionJoin represents joined party
	PartyActionJoin
)

// PartyActionResult represents a packet for party action results
type PartyActionResult struct {
	*packets.BasePacket
	PlayerID uint16
	ActionID PartyActionID
}

// NewPartyActionResult creates a new PartyActionResult packet
func NewPartyActionResult() *PartyActionResult {
	return &PartyActionResult{
		BasePacket: packets.NewPacket(interfaces.PartyActionResult, 254), // 254 is the unsigned byte equivalent of -2
	}
}

// Type returns the packet type
func (p *PartyActionResult) Type() interfaces.PacketType {
	return interfaces.PartyActionResult
}

// Read reads the packet data from a PacketReader
func (p *PartyActionResult) Read(r *packets.PacketReader) error {
	var err error
	p.PlayerID, err = r.ReadUInt16()
	if err != nil {
		return err
	}
	actionID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.ActionID = PartyActionID(actionID)
	return nil
}

// Write writes the packet data to a PacketWriter
func (p *PartyActionResult) Write(w *packets.PacketWriter) error {
	if err := w.WriteUInt16(p.PlayerID); err != nil {
		return err
	}
	return w.WriteByte(byte(p.ActionID))
}
