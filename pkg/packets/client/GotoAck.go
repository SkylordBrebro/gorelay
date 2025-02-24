package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// GotoAck represents a packet for acknowledging goto commands
type GotoAck struct {
	*packets.BasePacket
	Time    int32
	Unknown bool // UnknownFCOJFJPBJFA in original
}

// NewGotoAck creates a new GotoAck packet
func NewGotoAck() *GotoAck {
	return &GotoAck{
		BasePacket: packets.NewPacket(interfaces.GotoAck, byte(interfaces.GotoAck)),
	}
}

// Type returns the packet type
func (p *GotoAck) Type() interfaces.PacketType {
	return interfaces.GotoAck
}

// Read reads the packet data from a PacketReader
func (p *GotoAck) Read(r *packets.PacketReader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Unknown, err = r.ReadBool()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *GotoAck) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	return w.WriteBool(p.Unknown)
}
