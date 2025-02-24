package server

import (
	"gorelay/pkg/packets/interfaces"
)

// PartyMemberAdded represents the server packet for party member addition
type PartyMemberAdded struct {
	PlayerId uint16
	Name     string
	ClassId  uint16
	SkinId   uint16
}

// Type returns the packet type for PartyMemberAdded
func (p *PartyMemberAdded) Type() interfaces.PacketType {
	return interfaces.PartyMemberAdded
}

// Read reads the packet data from the provided reader
func (p *PartyMemberAdded) Read(r interfaces.Reader) error {
	var err error

	// Read PlayerId
	p.PlayerId, err = r.ReadUInt16()
	if err != nil {
		return err
	}

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

	return nil
}

// Write writes the packet data to the provided writer
func (p *PartyMemberAdded) Write(w interfaces.Writer) error {
	var err error

	// Write PlayerId
	err = w.WriteUInt16(p.PlayerId)
	if err != nil {
		return err
	}

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

	return nil
}
