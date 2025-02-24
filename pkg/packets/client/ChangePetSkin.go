package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// ChangePetSkin represents a packet for changing a pet's skin
type ChangePetSkin struct {
	*packets.BasePacket
	PetID    int32
	SkinType int32
	Currency int32
}

// NewChangePetSkin creates a new ChangePetSkin packet
func NewChangePetSkin() *ChangePetSkin {
	return &ChangePetSkin{
		BasePacket: packets.NewPacket(interfaces.ChangePetSkin, byte(interfaces.ChangePetSkin)),
	}
}

// Type returns the packet type
func (p *ChangePetSkin) Type() interfaces.PacketType {
	return interfaces.ChangePetSkin
}

// Read reads the packet data from a PacketReader
func (p *ChangePetSkin) Read(r *packets.PacketReader) error {
	var err error
	p.PetID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.SkinType, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Currency, err = r.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *ChangePetSkin) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.PetID); err != nil {
		return err
	}
	if err := w.WriteInt32(p.SkinType); err != nil {
		return err
	}
	return w.WriteInt32(p.Currency)
}
