package client

import (
	"gorelay/pkg/packets/interfaces"
)

// UsePortal represents a client-side use portal packet
type UsePortal struct {
	ObjectId int32
}

// Type returns the packet type for UsePortal
func (p *UsePortal) Type() interfaces.PacketType {
	return interfaces.UsePortal
}

// Read reads the packet data from the given reader
func (p *UsePortal) Read(r interfaces.Reader) error {
	var err error
	p.ObjectId, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the given writer
func (p *UsePortal) Write(w interfaces.Writer) error {
	return w.WriteInt32(p.ObjectId)
}
