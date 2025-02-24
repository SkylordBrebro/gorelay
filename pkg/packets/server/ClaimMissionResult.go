package server

import (
	"gorelay/pkg/packets/interfaces"
)

// ClaimMissionResult represents a server-side claim mission result packet
type ClaimMissionResult struct {
	MissionType byte
	Success     bool
	Message     string
}

// Type returns the packet type for ClaimMissionResult
func (p *ClaimMissionResult) Type() interfaces.PacketType {
	return interfaces.ClaimMissionResult
}

// ID returns the packet ID
func (p *ClaimMissionResult) ID() int32 {
	return int32(interfaces.ClaimMissionResult)
}

// Read reads the packet data from the given reader
func (p *ClaimMissionResult) Read(r interfaces.Reader) error {
	var err error
	p.MissionType, err = r.ReadByte()
	if err != nil {
		return err
	}

	p.Success, err = r.ReadBool()
	if err != nil {
		return err
	}

	p.Message, err = r.ReadString()
	return err
}

// Write writes the packet data to the given writer
func (p *ClaimMissionResult) Write(w interfaces.Writer) error {
	if err := w.WriteByte(p.MissionType); err != nil {
		return err
	}

	if err := w.WriteBool(p.Success); err != nil {
		return err
	}

	return w.WriteString(p.Message)
}

// String returns a string representation of the packet
func (p *ClaimMissionResult) String() string {
	return "ClaimMissionResult"
}

// HasNulls checks if any fields in the packet are null
func (p *ClaimMissionResult) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *ClaimMissionResult) Structure() string {
	return "ClaimMissionResult"
}
