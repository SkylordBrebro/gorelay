package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// PlayerText represents a packet for player text messages
type PlayerText struct {
	*packets.BasePacket
	Text string
}

// NewPlayerText creates a new PlayerText packet
func NewPlayerText() *PlayerText {
	return &PlayerText{
		BasePacket: packets.NewPacket(interfaces.PlayerText, byte(interfaces.PlayerText)),
	}
}

// Type returns the packet type
func (p *PlayerText) Type() interfaces.PacketType {
	return interfaces.PlayerText
}

// Read reads the packet data from a PacketReader
func (p *PlayerText) Read(r *packets.PacketReader) error {
	var err error
	p.Text, err = r.ReadString()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *PlayerText) Write(w *packets.PacketWriter) error {
	return w.WriteString(p.Text)
}
