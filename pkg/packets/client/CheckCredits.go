package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// CheckCredits represents a packet for checking credits
type CheckCredits struct {
	*packets.BasePacket
}

// NewCheckCredits creates a new CheckCredits packet
func NewCheckCredits() *CheckCredits {
	return &CheckCredits{
		BasePacket: packets.NewPacket(interfaces.CheckCredits, byte(interfaces.CheckCredits)),
	}
}

// Type returns the packet type
func (p *CheckCredits) Type() interfaces.PacketType {
	return interfaces.CheckCredits
}

// Read reads the packet data from a PacketReader
func (p *CheckCredits) Read(r *packets.PacketReader) error {
	return nil
}

// Write writes the packet data to a PacketWriter
func (p *CheckCredits) Write(w *packets.PacketWriter) error {
	return nil
}
