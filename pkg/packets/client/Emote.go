package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// Emote represents a packet for performing an emote
type Emote struct {
	*packets.BasePacket
	EmoteID      int32
	Time         int32
	UnknownBool0 bool
}

// NewEmote creates a new Emote packet
func NewEmote() *Emote {
	return &Emote{
		BasePacket: packets.NewPacket(interfaces.Emote, byte(interfaces.Emote)),
	}
}

// Type returns the packet type
func (p *Emote) Type() interfaces.PacketType {
	return interfaces.Emote
}

// Read reads the packet data from a PacketReader
func (p *Emote) Read(r *packets.PacketReader) error {
	var err error
	p.EmoteID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.UnknownBool0, err = r.ReadBool()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *Emote) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.EmoteID); err != nil {
		return err
	}
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	return w.WriteBool(p.UnknownBool0)
}
