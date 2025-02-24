package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// ChangeTrade represents a packet for modifying a trade
type ChangeTrade struct {
	*packets.BasePacket
	Offers []bool
}

// NewChangeTrade creates a new ChangeTrade packet
func NewChangeTrade() *ChangeTrade {
	return &ChangeTrade{
		BasePacket: packets.NewPacket(interfaces.ChangeTrade, byte(interfaces.ChangeTrade)),
		Offers:     make([]bool, 0),
	}
}

// Type returns the packet type
func (p *ChangeTrade) Type() interfaces.PacketType {
	return interfaces.ChangeTrade
}

// Read reads the packet data from a PacketReader
func (p *ChangeTrade) Read(r *packets.PacketReader) error {
	length, err := r.ReadInt16()
	if err != nil {
		return err
	}

	p.Offers = make([]bool, length)
	for i := 0; i < int(length); i++ {
		p.Offers[i], err = r.ReadBool()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to a PacketWriter
func (p *ChangeTrade) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt16(int16(len(p.Offers))); err != nil {
		return err
	}

	for _, offer := range p.Offers {
		if err := w.WriteBool(offer); err != nil {
			return err
		}
	}

	return nil
}
