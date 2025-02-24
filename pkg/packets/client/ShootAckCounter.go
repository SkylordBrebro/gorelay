package client

import (
	"gorelay/pkg/packets/interfaces"
)

// ShootAckCounter represents a client-side shoot acknowledgment counter packet
type ShootAckCounter struct {
	Time   int32
	Amount int16
}

// Type returns the packet type for ShootAckCounter
func (p *ShootAckCounter) Type() interfaces.PacketType {
	return interfaces.ShootAckCounter
}

// Read reads the packet data from the given reader
func (p *ShootAckCounter) Read(r interfaces.Reader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Amount, err = r.ReadInt16()
	return err
}

// Write writes the packet data to the given writer
func (p *ShootAckCounter) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	return w.WriteInt16(p.Amount)
}
