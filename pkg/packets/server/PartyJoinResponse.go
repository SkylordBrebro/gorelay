package server

import (
	"gorelay/pkg/packets/interfaces"
)

// PartyJoinResponse represents the server packet for party join responses
type PartyJoinResponse struct {
	ServerIpHost string
}

// Type returns the packet type for PartyJoinResponse
func (p *PartyJoinResponse) Type() interfaces.PacketType {
	return interfaces.PartyJoinResponse
}

// Read reads the packet data from the provided reader
func (p *PartyJoinResponse) Read(r interfaces.Reader) error {
	var err error
	p.ServerIpHost, err = r.ReadString()
	return err
}

// Write writes the packet data to the provided writer
func (p *PartyJoinResponse) Write(w interfaces.Writer) error {
	return w.WriteString(p.ServerIpHost)
}

func (p *PartyJoinResponse) ID() int32 {
	return int32(interfaces.PartyJoinResponse)
}