package client

import (
	"gorelay/pkg/packets/interfaces"
)

// UpdateAck represents a client-side update acknowledgment packet
type UpdateAck struct{}

// Default is the default instance of UpdateAck
var Default = &UpdateAck{}

// Type returns the packet type for UpdateAck
func (p *UpdateAck) Type() interfaces.PacketType {
	return interfaces.UpdateAck
}

// Read reads the packet data from the given reader
func (p *UpdateAck) Read(r interfaces.Reader) error {
	return nil
}

// Write writes the packet data to the given writer
func (p *UpdateAck) Write(w interfaces.Writer) error {
	return nil
}
