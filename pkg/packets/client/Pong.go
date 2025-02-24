package client

import (
	"gorelay/pkg/packets/interfaces"
)

// Pong represents a client-side pong packet
type Pong struct {
	Serial int32
	Time   int32
}

// Type returns the packet type for Pong
func (p *Pong) Type() interfaces.PacketType {
	return interfaces.Pong
}

// Read reads the packet data from the given reader
func (p *Pong) Read(r interfaces.Reader) error {
	var err error
	p.Serial, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Time, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the given writer
func (p *Pong) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.Serial); err != nil {
		return err
	}
	return w.WriteInt32(p.Time)
}
