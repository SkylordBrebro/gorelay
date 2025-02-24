package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// BoostBPMilestone represents a battle pass milestone boost packet
type BoostBPMilestone struct {
	*packets.BasePacket
	MilestoneID byte
}

// NewBoostBPMilestone creates a new BoostBPMilestone packet
func NewBoostBPMilestone() *BoostBPMilestone {
	return &BoostBPMilestone{
		BasePacket: packets.NewPacket(interfaces.BoostBPMilestone, byte(interfaces.BoostBPMilestone)),
	}
}

// Type returns the packet type
func (p *BoostBPMilestone) Type() interfaces.PacketType {
	return interfaces.BoostBPMilestone
}

// Read reads the packet data from the reader
func (p *BoostBPMilestone) Read(r *packets.PacketReader) error {
	var err error
	p.MilestoneID, err = r.ReadByte()
	return err
}

// Write writes the packet data to the writer
func (p *BoostBPMilestone) Write(w *packets.PacketWriter) error {
	return w.WriteByte(p.MilestoneID)
}
