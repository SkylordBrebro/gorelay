package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Ping represents the server packet for ping
type Ping struct {
	Serial int32
}

// Type returns the packet type for Ping
func (p *Ping) Type() interfaces.PacketType {
	return interfaces.Ping
}

// Read reads the packet data from the provided reader
func (p *Ping) Read(r interfaces.Reader) error {
	var err error
	p.Serial, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the provided writer
func (p *Ping) Write(w interfaces.Writer) error {
	return w.WriteInt32(p.Serial)
}

func (p *Ping) ID() int32 {
	return int32(interfaces.Ping)
}