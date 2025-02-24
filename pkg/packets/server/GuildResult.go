package server

import (
	"gorelay/pkg/packets/interfaces"
)

// GuildResult represents the server packet for guild operation results
type GuildResult struct {
	Success   bool
	ErrorText string
}

// Type returns the packet type for GuildResult
func (p *GuildResult) Type() interfaces.PacketType {
	return interfaces.GuildResult
}

// Read reads the packet data from the provided reader
func (p *GuildResult) Read(r interfaces.Reader) error {
	var err error

	// Read Success
	p.Success, err = r.ReadBool()
	if err != nil {
		return err
	}

	// Read ErrorText
	p.ErrorText, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *GuildResult) Write(w interfaces.Writer) error {
	var err error

	// Write Success
	err = w.WriteBool(p.Success)
	if err != nil {
		return err
	}

	// Write ErrorText
	err = w.WriteString(p.ErrorText)
	if err != nil {
		return err
	}

	return nil
}
