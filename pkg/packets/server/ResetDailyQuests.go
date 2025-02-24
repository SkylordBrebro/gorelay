package server

import (
	"gorelay/pkg/packets/interfaces"
)

// ResetDailyQuests represents the server packet for resetting daily quests
type ResetDailyQuests struct {
	// This packet has no fields
}

// Type returns the packet type for ResetDailyQuests
func (p *ResetDailyQuests) Type() interfaces.PacketType {
	return interfaces.ResetDailyQuests
}

// Read reads the packet data from the provided reader
func (p *ResetDailyQuests) Read(r interfaces.Reader) error {
	// No data to read
	return nil
}

// Write writes the packet data to the provided writer
func (p *ResetDailyQuests) Write(w interfaces.Writer) error {
	// No data to write
	return nil
}
