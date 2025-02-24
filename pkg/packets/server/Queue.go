package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Queue represents the server packet for queue information
type Queue struct {
	CurrentPosition uint16
	MaxPosition     uint16
}

// Type returns the packet type for Queue
func (p *Queue) Type() interfaces.PacketType {
	return interfaces.Queue
}

// Read reads the packet data from the provided reader
func (p *Queue) Read(r interfaces.Reader) error {
	var err error

	// Read CurrentPosition
	p.CurrentPosition, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	// Read MaxPosition
	p.MaxPosition, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *Queue) Write(w interfaces.Writer) error {
	var err error

	// Write CurrentPosition
	err = w.WriteUInt16(p.CurrentPosition)
	if err != nil {
		return err
	}

	// Write MaxPosition
	err = w.WriteUInt16(p.MaxPosition)
	if err != nil {
		return err
	}

	return nil
}
