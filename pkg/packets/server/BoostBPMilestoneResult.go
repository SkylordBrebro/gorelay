package server

import (
	"gorelay/pkg/packets/interfaces"
)

// BoostBPMilestoneResult represents a server-side boost BP milestone result packet
type BoostBPMilestoneResult struct {
	Success bool
}

// Type returns the packet type for BoostBPMilestoneResult
func (p *BoostBPMilestoneResult) Type() interfaces.PacketType {
	return interfaces.BoostBPMilestoneResult
}

// ID returns the packet ID
func (p *BoostBPMilestoneResult) ID() int32 {
	return int32(interfaces.BoostBPMilestoneResult)
}

// Read reads the packet data from the given reader
func (p *BoostBPMilestoneResult) Read(r interfaces.Reader) error {
	var err error
	p.Success, err = r.ReadBool()
	return err
}

// Write writes the packet data to the given writer
func (p *BoostBPMilestoneResult) Write(w interfaces.Writer) error {
	return w.WriteBool(p.Success)
}

// String returns a string representation of the packet
func (p *BoostBPMilestoneResult) String() string {
	return "BoostBPMilestoneResult"
}

// HasNulls checks if any fields in the packet are null
func (p *BoostBPMilestoneResult) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *BoostBPMilestoneResult) Structure() string {
	return "BoostBPMilestoneResult"
}
