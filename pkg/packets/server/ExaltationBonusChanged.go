package server

import (
	"gorelay/pkg/packets/interfaces"
)

// ExaltationBonusChanged represents the server packet for exaltation bonus changes
type ExaltationBonusChanged struct {
	// This packet appears to have no fields in the original C# implementation
}

// Type returns the packet type for ExaltationBonusChanged
func (p *ExaltationBonusChanged) Type() interfaces.PacketType {
	return interfaces.ExaltationBonusChanged
}

// Read reads the packet data from the provided reader
func (p *ExaltationBonusChanged) Read(r interfaces.Reader) error {
	// No fields to read
	return nil
}

// Write writes the packet data to the provided writer
func (p *ExaltationBonusChanged) Write(w interfaces.Writer) error {
	// No fields to write
	return nil
}

func (p *ExaltationBonusChanged) ID() int32 {
	return int32(interfaces.ExaltationBonusChanged)
}