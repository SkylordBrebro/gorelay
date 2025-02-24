package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// ClaimMission represents a packet for claiming a mission
type ClaimMission struct {
	*packets.BasePacket
	MissionID   int32  // CIGJNLDGNAF in original
	MissionType byte   // BGJKHDCELPO in original
	Category    byte   // MAOGDBIHOOB in original
	SubCategory uint16 // JOCFHNELMFM in original
}

// NewClaimMission creates a new ClaimMission packet
func NewClaimMission() *ClaimMission {
	return &ClaimMission{
		BasePacket: packets.NewPacket(interfaces.ClaimMission, byte(interfaces.ClaimMission)),
	}
}

// Type returns the packet type
func (p *ClaimMission) Type() interfaces.PacketType {
	return interfaces.ClaimMission
}

// Read reads the packet data from a PacketReader
func (p *ClaimMission) Read(r *packets.PacketReader) error {
	var err error
	p.MissionID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.MissionType, err = r.ReadByte()
	if err != nil {
		return err
	}
	p.Category, err = r.ReadByte()
	if err != nil {
		return err
	}
	p.SubCategory, err = r.ReadUInt16()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *ClaimMission) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.MissionID); err != nil {
		return err
	}
	if err := w.WriteByte(p.MissionType); err != nil {
		return err
	}
	if err := w.WriteByte(p.Category); err != nil {
		return err
	}
	return w.WriteUInt16(p.SubCategory)
}
