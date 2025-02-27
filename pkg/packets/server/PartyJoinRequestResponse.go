package server

import (
	"gorelay/pkg/packets/interfaces"
)

// InviteState represents different party invite states
type InviteState byte

// Invite state constants
const (
	InviteStateNone        InviteState = 0
	InviteStatePending     InviteState = 1
	InviteStateCancelled   InviteState = 2
	InviteStateAccepted    InviteState = 3
	InviteStateDeclined    InviteState = 4
	InviteStatePartyFull   InviteState = 5
	InviteStateBlacklisted InviteState = 6
)

// PartyJoinRequestResponse represents the server packet for party join request responses
type PartyJoinRequestResponse struct {
	Name    string
	ClassId uint16
	SkinId  uint16
	State   InviteState
}

// Type returns the packet type for PartyJoinRequestResponse
func (p *PartyJoinRequestResponse) Type() interfaces.PacketType {
	return interfaces.PartyJoinRequestResponse
}

// Read reads the packet data from the provided reader
func (p *PartyJoinRequestResponse) Read(r interfaces.Reader) error {
	var err error

	// Read Name
	p.Name, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read ClassId
	p.ClassId, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	// Read SkinId
	p.SkinId, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	// Read State
	stateValue, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.State = InviteState(stateValue)

	return nil
}

// Write writes the packet data to the provided writer
func (p *PartyJoinRequestResponse) Write(w interfaces.Writer) error {
	var err error

	// Write Name
	err = w.WriteString(p.Name)
	if err != nil {
		return err
	}

	// Write ClassId
	err = w.WriteUInt16(p.ClassId)
	if err != nil {
		return err
	}

	// Write SkinId
	err = w.WriteUInt16(p.SkinId)
	if err != nil {
		return err
	}

	// Write State
	err = w.WriteByte(byte(p.State))
	if err != nil {
		return err
	}

	return nil
}

func (p *PartyJoinRequestResponse) ID() int32 {
	return int32(interfaces.PartyJoinRequestResponse)
}