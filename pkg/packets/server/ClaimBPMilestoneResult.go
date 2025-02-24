package server

import (
	"gorelay/pkg/packets/interfaces"
)

// ClaimBPMilestoneResult represents a server-side claim BP milestone result packet
type ClaimBPMilestoneResult struct {
	Success bool
}

// Type returns the packet type for ClaimBPMilestoneResult
func (p *ClaimBPMilestoneResult) Type() interfaces.PacketType {
	return interfaces.ClaimBPMilestoneResult
}

// ID returns the packet ID
func (p *ClaimBPMilestoneResult) ID() int32 {
	return int32(interfaces.ClaimBPMilestoneResult)
}

// Read reads the packet data from the given reader
func (p *ClaimBPMilestoneResult) Read(r interfaces.Reader) error {
	var err error
	p.Success, err = r.ReadBool()
	return err
}

// Write writes the packet data to the given writer
func (p *ClaimBPMilestoneResult) Write(w interfaces.Writer) error {
	return w.WriteBool(p.Success)
}

// String returns a string representation of the packet
func (p *ClaimBPMilestoneResult) String() string {
	return "ClaimBPMilestoneResult"
}

// HasNulls checks if any fields in the packet are null
func (p *ClaimBPMilestoneResult) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *ClaimBPMilestoneResult) Structure() string {
	return "ClaimBPMilestoneResult"
}
