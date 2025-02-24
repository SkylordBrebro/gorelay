package server

import (
	"gorelay/pkg/packets/interfaces"
)

// CrucibleResult represents the server packet for crucible operation results
type CrucibleResult struct {
	Success bool
}

// Type returns the packet type for CrucibleResult
func (p *CrucibleResult) Type() interfaces.PacketType {
	return interfaces.CrucibleResult
}

// Read reads the packet data from the provided reader
func (p *CrucibleResult) Read(r interfaces.Reader) error {
	var err error
	p.Success, err = r.ReadBool()
	return err
}

// Write writes the packet data to the provided writer
func (p *CrucibleResult) Write(w interfaces.Writer) error {
	return w.WriteBool(p.Success)
}
