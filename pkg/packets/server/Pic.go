package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Pic represents the server packet for picture data
type Pic struct {
	Data []byte
}

// Type returns the packet type for Pic
func (p *Pic) Type() interfaces.PacketType {
	return interfaces.Pic
}

// Read reads the packet data from the provided reader
func (p *Pic) Read(r interfaces.Reader) error {
	var err error

	// Read all remaining bytes in the stream (minus 5 bytes for header)
	remainingBytes := int(r.RemainingBytes())
	if remainingBytes > 0 {
		p.Data, err = r.ReadBytes(remainingBytes)
	} else {
		p.Data = []byte{}
	}

	return err
}

// Write writes the packet data to the provided writer
func (p *Pic) Write(w interfaces.Writer) error {
	return w.WriteBytes(p.Data)
}

func (p *Pic) ID() int32 {
	return int32(interfaces.Pic)
}