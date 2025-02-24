package server

import (
	"gorelay/pkg/packets/interfaces"
)

// TradeRequested represents the server packet for trade requests
type TradeRequested struct {
	Name string
}

// Type returns the packet type for TradeRequested
func (p *TradeRequested) Type() interfaces.PacketType {
	return interfaces.TradeRequested
}

// Read reads the packet data from the provided reader
func (p *TradeRequested) Read(r interfaces.Reader) error {
	var err error
	p.Name, err = r.ReadString()
	return err
}

// Write writes the packet data to the provided writer
func (p *TradeRequested) Write(w interfaces.Writer) error {
	return w.WriteString(p.Name)
}
