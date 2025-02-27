package server

import (
	"gorelay/pkg/packets/interfaces"
)

// Failure represents the server packet for failure notifications
type Failure struct {
	ErrorId      int32
	ErrorMessage string
}

// Type returns the packet type for Failure
func (p *Failure) Type() interfaces.PacketType {
	return interfaces.Failure
}

// Read reads the packet data from the provided reader
func (p *Failure) Read(r interfaces.Reader) error {
	var err error

	p.ErrorId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.ErrorMessage, err = r.ReadString()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *Failure) Write(w interfaces.Writer) error {
	var err error

	err = w.WriteInt32(p.ErrorId)
	if err != nil {
		return err
	}

	err = w.WriteString(p.ErrorMessage)
	if err != nil {
		return err
	}

	return nil
}

func (p *Failure) ID() int32 {
	return int32(interfaces.Failure)
}