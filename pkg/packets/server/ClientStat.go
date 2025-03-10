package server

import (
	"gorelay/pkg/packets/interfaces"
)

// ClientStat represents a client statistics packet from the server
type ClientStat struct {
	Name  string
	Value int32
}

// Type returns the packet type for ClientStat
func (p *ClientStat) Type() interfaces.PacketType {
	return interfaces.ClientStat
}

// Read reads the packet data from the provided reader
func (p *ClientStat) Read(r interfaces.Reader) error {
	var err error

	p.Name, err = r.ReadString()
	if err != nil {
		return err
	}

	p.Value, err = r.ReadInt32()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *ClientStat) Write(w interfaces.Writer) error {
	var err error

	err = w.WriteString(p.Name)
	if err != nil {
		return err
	}

	err = w.WriteInt32(p.Value)
	if err != nil {
		return err
	}

	return nil
}

func (p *ClientStat) ID() int32 {
	return int32(interfaces.ClientStat)
}

// String returns a string representation of the packet
func (p *ClientStat) String() string {
	return "ClientStat"
}

// HasNulls checks if any fields in the packet are null
func (p *ClientStat) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *ClientStat) Structure() string {
	return "ClientStat"
}
