package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// PartyPrivacy represents party privacy settings
type PartyPrivacy byte

// PartyCreate represents a packet for creating a party
type PartyCreate struct {
	*packets.BasePacket
	Description   string
	PowerLevelMin uint16
	PartySizeMax  byte
	Activity      byte
	StatsMaxedMin byte
	Privacy       PartyPrivacy
	ServerIndex   byte
}

// NewPartyCreate creates a new PartyCreate packet
func NewPartyCreate() *PartyCreate {
	return &PartyCreate{
		BasePacket: packets.NewPacket(interfaces.PartyAction, 251), // 251 is the unsigned byte equivalent of -5
	}
}

// Type returns the packet type
func (p *PartyCreate) Type() interfaces.PacketType {
	return interfaces.PartyAction
}

// Read reads the packet data from a PacketReader
func (p *PartyCreate) Read(r *packets.PacketReader) error {
	var err error
	p.Description, err = r.ReadString()
	if err != nil {
		return err
	}
	p.PowerLevelMin, err = r.ReadUInt16()
	if err != nil {
		return err
	}
	p.PartySizeMax, err = r.ReadByte()
	if err != nil {
		return err
	}
	p.Activity, err = r.ReadByte()
	if err != nil {
		return err
	}
	p.StatsMaxedMin, err = r.ReadByte()
	if err != nil {
		return err
	}
	privacyByte, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Privacy = PartyPrivacy(privacyByte)
	p.ServerIndex, err = r.ReadByte()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *PartyCreate) Write(w *packets.PacketWriter) error {
	if err := w.WriteString(p.Description); err != nil {
		return err
	}
	if err := w.WriteUInt16(p.PowerLevelMin); err != nil {
		return err
	}
	if err := w.WriteByte(p.PartySizeMax); err != nil {
		return err
	}
	if err := w.WriteByte(p.Activity); err != nil {
		return err
	}
	if err := w.WriteByte(p.StatsMaxedMin); err != nil {
		return err
	}
	if err := w.WriteByte(byte(p.Privacy)); err != nil {
		return err
	}
	return w.WriteByte(p.ServerIndex)
}
