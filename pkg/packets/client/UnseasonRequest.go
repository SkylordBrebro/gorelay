package client

import (
	"gorelay/pkg/packets/interfaces"
)

// UnseasonRequest represents a client-side unseason request packet
type UnseasonRequest struct{}

// Type returns the packet type for UnseasonRequest
func (p *UnseasonRequest) Type() interfaces.PacketType {
	return interfaces.UnseasonRequest
}

// Read reads the packet data from the given reader
func (p *UnseasonRequest) Read(r interfaces.Reader) error {
	return nil
}

// Write writes the packet data to the given writer
func (p *UnseasonRequest) Write(w interfaces.Writer) error {
	return nil
}
