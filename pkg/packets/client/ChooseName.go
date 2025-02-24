package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// ChooseName represents a packet for choosing a character name
type ChooseName struct {
	*packets.BasePacket
	Name string
}

// NewChooseName creates a new ChooseName packet
func NewChooseName() *ChooseName {
	return &ChooseName{
		BasePacket: packets.NewPacket(interfaces.ChooseName, byte(interfaces.ChooseName)),
	}
}

// Type returns the packet type
func (p *ChooseName) Type() interfaces.PacketType {
	return interfaces.ChooseName
}

// Read reads the packet data from a PacketReader
func (p *ChooseName) Read(r *packets.PacketReader) error {
	var err error
	p.Name, err = r.ReadString()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *ChooseName) Write(w *packets.PacketWriter) error {
	return w.WriteString(p.Name)
}
