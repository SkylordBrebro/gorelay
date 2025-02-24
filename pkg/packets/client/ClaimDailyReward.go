package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// ClaimDailyReward represents a packet for claiming daily rewards
type ClaimDailyReward struct {
	*packets.BasePacket
	ClaimKey  string
	ClaimType string
}

// NewClaimDailyReward creates a new ClaimDailyReward packet
func NewClaimDailyReward() *ClaimDailyReward {
	return &ClaimDailyReward{
		BasePacket: packets.NewPacket(interfaces.ClaimDailyReward, byte(interfaces.ClaimDailyReward)),
	}
}

// Type returns the packet type
func (p *ClaimDailyReward) Type() interfaces.PacketType {
	return interfaces.ClaimDailyReward
}

// Read reads the packet data from a PacketReader
func (p *ClaimDailyReward) Read(r *packets.PacketReader) error {
	var err error
	p.ClaimKey, err = r.ReadString()
	if err != nil {
		return err
	}
	p.ClaimType, err = r.ReadString()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *ClaimDailyReward) Write(w *packets.PacketWriter) error {
	if err := w.WriteString(p.ClaimKey); err != nil {
		return err
	}
	return w.WriteString(p.ClaimType)
}
