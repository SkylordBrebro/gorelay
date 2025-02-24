package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// ClaimBPMilestone represents a packet for claiming a battle pass milestone
type ClaimBPMilestone struct {
	*packets.BasePacket
	RewardID int8 // Using int8 instead of sbyte as Go doesn't have sbyte
}

// NewClaimBPMilestone creates a new ClaimBPMilestone packet
func NewClaimBPMilestone() *ClaimBPMilestone {
	return &ClaimBPMilestone{
		BasePacket: packets.NewPacket(interfaces.ClaimBPMilestone, byte(interfaces.ClaimBPMilestone)),
	}
}

// Type returns the packet type
func (p *ClaimBPMilestone) Type() interfaces.PacketType {
	return interfaces.ClaimBPMilestone
}

// Read reads the packet data from a PacketReader
func (p *ClaimBPMilestone) Read(r *packets.PacketReader) error {
	rewardID, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.RewardID = int8(rewardID)
	return nil
}

// Write writes the packet data to a PacketWriter
func (p *ClaimBPMilestone) Write(w *packets.PacketWriter) error {
	return w.WriteByte(byte(p.RewardID))
}
