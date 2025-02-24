package client

import (
	"gorelay/pkg/packets/interfaces"
)

// Teleport represents a client-side teleport packet
type Teleport struct {
	ObjectId int32
	Name     string
}

// Type returns the packet type for Teleport
func (p *Teleport) Type() interfaces.PacketType {
	return interfaces.Teleport
}

// Read reads the packet data from the given reader
func (p *Teleport) Read(r interfaces.Reader) error {
	var err error
	p.ObjectId, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Name, err = r.ReadString()
	return err
}

// Write writes the packet data to the given writer
func (p *Teleport) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.ObjectId); err != nil {
		return err
	}
	return w.WriteString(p.Name)
}
