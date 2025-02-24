package server

import (
	"gorelay/pkg/packets/interfaces"
)

// NameResult represents the server packet for name change results
type NameResult struct {
	Success   bool
	ErrorText string
}

// Type returns the packet type for NameResult
func (p *NameResult) Type() interfaces.PacketType {
	return interfaces.NameResult
}

// Read reads the packet data from the provided reader
func (p *NameResult) Read(r interfaces.Reader) error {
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
func (p *NameResult) Write(w interfaces.Writer) error {
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
