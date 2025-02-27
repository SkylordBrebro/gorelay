package server

import (
	"gorelay/pkg/packets/interfaces"
)

// InvitedToGuild represents the server packet for guild invitations
type InvitedToGuild struct {
	Name      string
	GuildName string
}

// Type returns the packet type for InvitedToGuild
func (p *InvitedToGuild) Type() interfaces.PacketType {
	return interfaces.InvitedToGuild
}

// Read reads the packet data from the provided reader
func (p *InvitedToGuild) Read(r interfaces.Reader) error {
	var err error

	// Read Name
	p.Name, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read GuildName
	p.GuildName, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *InvitedToGuild) Write(w interfaces.Writer) error {
	var err error

	// Write Name
	err = w.WriteString(p.Name)
	if err != nil {
		return err
	}

	// Write GuildName
	err = w.WriteString(p.GuildName)
	if err != nil {
		return err
	}

	return nil
}

func (p *InvitedToGuild) ID() int32 {
	return int32(interfaces.InvitedToGuild)
}