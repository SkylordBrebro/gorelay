package server

import (
	"gorelay/pkg/packets/interfaces"
)

// PartyActionId represents different party action types
type PartyActionId byte

// Party action type constants
const (
	PartyActionNone             PartyActionId = 0
	PartyActionFailed           PartyActionId = 1
	PartyActionKicked           PartyActionId = 2
	PartyActionKickNotFound     PartyActionId = 3
	PartyActionPromotedToLeader PartyActionId = 4
	PartyActionPromoteNotFound  PartyActionId = 5
	PartyActionLeftParty        PartyActionId = 6
	PartyActionJoin             PartyActionId = 7
)

// PartyAction represents the server packet for party actions
type PartyAction struct {
	PlayerId uint16
	ActionId PartyActionId
}

// Type returns the packet type for PartyAction
func (p *PartyAction) Type() interfaces.PacketType {
	return interfaces.PartyAction
}

// Read reads the packet data from the provided reader
func (p *PartyAction) Read(r interfaces.Reader) error {
	var err error

	// Read PlayerId
	p.PlayerId, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	// Read ActionId
	actionId, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.ActionId = PartyActionId(actionId)

	return nil
}

// Write writes the packet data to the provided writer
func (p *PartyAction) Write(w interfaces.Writer) error {
	var err error

	// Write PlayerId
	err = w.WriteUInt16(p.PlayerId)
	if err != nil {
		return err
	}

	// Write ActionId
	err = w.WriteByte(byte(p.ActionId))
	if err != nil {
		return err
	}

	return nil
}

func (p *PartyAction) ID() int32 {
	return int32(interfaces.PartyAction)
}