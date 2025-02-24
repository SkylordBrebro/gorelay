package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// Move represents a packet for movement
type Move struct {
	*packets.BasePacket
	TickID  int32
	Time    int32
	Records []*dataobjects.LocationRecord
}

// NewMove creates a new Move packet
func NewMove() *Move {
	return &Move{
		BasePacket: packets.NewPacket(interfaces.Move, byte(interfaces.Move)),
		Records:    make([]*dataobjects.LocationRecord, 0),
	}
}

// Type returns the packet type
func (p *Move) Type() interfaces.PacketType {
	return interfaces.Move
}

// Read reads the packet data from a PacketReader
func (p *Move) Read(r *packets.PacketReader) error {
	var err error
	p.TickID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}

	length, err := r.ReadInt16()
	if err != nil {
		return err
	}

	p.Records = make([]*dataobjects.LocationRecord, length)
	for i := 0; i < int(length); i++ {
		p.Records[i] = dataobjects.NewLocationRecord()
		if err := p.Records[i].Read(r); err != nil {
			return err
		}
	}
	return nil
}

// Write writes the packet data to a PacketWriter
func (p *Move) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.TickID); err != nil {
		return err
	}
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	if err := w.WriteInt16(int16(len(p.Records))); err != nil {
		return err
	}
	for _, record := range p.Records {
		if err := record.Write(w); err != nil {
			return err
		}
	}
	return nil
}
