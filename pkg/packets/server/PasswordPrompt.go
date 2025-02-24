package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Password prompt status constants
const (
	PasswordPromptSignIn             = 2
	PasswordPromptSendEmailAndSignIn = 3
	PasswordPromptRegister           = 4
)

// PasswordPrompt represents the server packet for password prompts
type PasswordPrompt struct {
	CleanPasswordStatus int32
}

// Type returns the packet type for PasswordPrompt
func (p *PasswordPrompt) Type() interfaces.PacketType {
	return interfaces.PasswordPrompt
}

// Read reads the packet data from the provided reader
func (p *PasswordPrompt) Read(r interfaces.Reader) error {
	var err error
	p.CleanPasswordStatus, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the provided writer
func (p *PasswordPrompt) Write(w interfaces.Writer) error {
	return w.WriteInt32(p.CleanPasswordStatus)
}
