package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// RedeemExaltationReward represents a packet for redeeming exaltation rewards
type RedeemExaltationReward struct {
	*packets.BasePacket
	ClassID int32
}

// NewRedeemExaltationReward creates a new RedeemExaltationReward packet
func NewRedeemExaltationReward() *RedeemExaltationReward {
	return &RedeemExaltationReward{
		BasePacket: packets.NewPacket(interfaces.RedeemExaltationReward, byte(interfaces.RedeemExaltationReward)),
	}
}

// Type returns the packet type
func (r *RedeemExaltationReward) Type() interfaces.PacketType {
	return interfaces.RedeemExaltationReward
}

// Read reads the packet data from a PacketReader
func (r *RedeemExaltationReward) Read(reader *packets.PacketReader) error {
	var err error
	r.ClassID, err = reader.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (r *RedeemExaltationReward) Write(writer *packets.PacketWriter) error {
	return writer.WriteInt32(r.ClassID)
}
