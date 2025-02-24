package server

import (
	"gorelay/pkg/packets/interfaces"
)

// IncomingPartyInvite represents the server packet for incoming party invitations
type IncomingPartyInvite struct {
	PartyId     uint32
	InviterName string
}

// Type returns the packet type for IncomingPartyInvite
func (p *IncomingPartyInvite) Type() interfaces.PacketType {
	return interfaces.IncomingPartyInvite
}

// Read reads the packet data from the provided reader
func (p *IncomingPartyInvite) Read(r interfaces.Reader) error {
	var err error

	// Read PartyId
	p.PartyId, err = r.ReadUInt32()
	if err != nil {
		return err
	}

	// Read InviterName
	p.InviterName, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *IncomingPartyInvite) Write(w interfaces.Writer) error {
	var err error

	// Write PartyId
	err = w.WriteUInt32(p.PartyId)
	if err != nil {
		return err
	}

	// Write InviterName
	err = w.WriteString(p.InviterName)
	if err != nil {
		return err
	}

	return nil
}
