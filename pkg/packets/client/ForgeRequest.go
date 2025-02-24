package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// ForgeRequest represents a packet for forge requests
type ForgeRequest struct {
	*packets.BasePacket
	ForgeTargetItem int32
	DismantleSlots  []*dataobjects.SlotObject
}

// NewForgeRequest creates a new ForgeRequest packet
func NewForgeRequest() *ForgeRequest {
	return &ForgeRequest{
		BasePacket:     packets.NewPacket(interfaces.ForgeRequest, byte(interfaces.ForgeRequest)),
		DismantleSlots: make([]*dataobjects.SlotObject, 0),
	}
}

// Type returns the packet type
func (p *ForgeRequest) Type() interfaces.PacketType {
	return interfaces.ForgeRequest
}

// Read reads the packet data from a PacketReader
func (p *ForgeRequest) Read(r *packets.PacketReader) error {
	var err error
	p.ForgeTargetItem, err = r.ReadInt32()
	if err != nil {
		return err
	}

	length, err := r.ReadInt32()
	if err != nil {
		return err
	}

	p.DismantleSlots = make([]*dataobjects.SlotObject, length)
	for i := 0; i < int(length); i++ {
		p.DismantleSlots[i] = dataobjects.NewSlotObject()
		if err := p.DismantleSlots[i].Read(r); err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to a PacketWriter
func (p *ForgeRequest) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.ForgeTargetItem); err != nil {
		return err
	}

	if err := w.WriteInt32(int32(len(p.DismantleSlots))); err != nil {
		return err
	}

	for _, slot := range p.DismantleSlots {
		if err := slot.Write(w); err != nil {
			return err
		}
	}

	return nil
}
