package server

import (
	"gorelay/pkg/packets/interfaces"
)

// HeroLeft represents the server packet for hero left notification
type HeroLeft struct {
	// This packet appears to have no fields in the original C# implementation
}

// Type returns the packet type for HeroLeft
func (p *HeroLeft) Type() interfaces.PacketType {
	return interfaces.HeroLeft
}

// Read reads the packet data from the provided reader
func (p *HeroLeft) Read(r interfaces.Reader) error {
	// No fields to read
	return nil
}

// Write writes the packet data to the provided writer
func (p *HeroLeft) Write(w interfaces.Writer) error {
	// No fields to write
	return nil
}

func (p *HeroLeft) ID() int32 {
	return int32(interfaces.HeroLeft)
}