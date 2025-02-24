package server

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// PlaySound represents a server packet for playing sounds
type PlaySound struct {
	*packets.BasePacket
	OwnerId int32
	SoundId byte
}

// NewPlaySound creates a new PlaySound packet
func NewPlaySound() *PlaySound {
	return &PlaySound{
		BasePacket: packets.NewPacket(interfaces.PlaySound, byte(interfaces.PlaySound)),
	}
}

// Type returns the packet type
func (p *PlaySound) Type() interfaces.PacketType {
	return interfaces.PlaySound
}

// Read reads the packet data from the reader
func (p *PlaySound) Read(r interfaces.Reader) error {
	var err error
	p.OwnerId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.SoundId, err = r.ReadByte()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the writer
func (p *PlaySound) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.OwnerId); err != nil {
		return err
	}

	if err := w.WriteByte(p.SoundId); err != nil {
		return err
	}

	return nil
}
