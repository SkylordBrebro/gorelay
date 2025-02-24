package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// PetUpgradeRequest represents a packet for requesting a pet upgrade
type PetUpgradeRequest struct {
	*packets.BasePacket
	PetTransType     byte
	PetID1           int32
	PetID2           int32
	ObjectID         int32
	ObjectSlot       int32
	PaymentTransType byte
}

// NewPetUpgradeRequest creates a new PetUpgradeRequest packet
func NewPetUpgradeRequest() *PetUpgradeRequest {
	return &PetUpgradeRequest{
		BasePacket: packets.NewPacket(interfaces.PetUpgradeRequest, byte(interfaces.PetUpgradeRequest)),
	}
}

// Type returns the packet type
func (p *PetUpgradeRequest) Type() interfaces.PacketType {
	return interfaces.PetUpgradeRequest
}

// Read reads the packet data from a PacketReader
func (p *PetUpgradeRequest) Read(r *packets.PacketReader) error {
	var err error
	p.PetTransType, err = r.ReadByte()
	if err != nil {
		return err
	}
	p.PetID1, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.PetID2, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.ObjectID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.ObjectSlot, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.PaymentTransType, err = r.ReadByte()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *PetUpgradeRequest) Write(w *packets.PacketWriter) error {
	if err := w.WriteByte(p.PetTransType); err != nil {
		return err
	}
	if err := w.WriteInt32(p.PetID1); err != nil {
		return err
	}
	if err := w.WriteInt32(p.PetID2); err != nil {
		return err
	}
	if err := w.WriteInt32(p.ObjectID); err != nil {
		return err
	}
	if err := w.WriteInt32(p.ObjectSlot); err != nil {
		return err
	}
	return w.WriteByte(p.PaymentTransType)
}
