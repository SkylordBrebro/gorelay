package server

import (
	"gorelay/pkg/packets/interfaces"
)

// CreateSuccess represents a server-side create success packet
type CreateSuccess struct {
	ObjectId int32
	CharId   int32
	Stats    string
}

// Type returns the packet type for CreateSuccess
func (p *CreateSuccess) Type() interfaces.PacketType {
	return interfaces.CreateSuccess
}

// ID returns the packet ID
func (p *CreateSuccess) ID() int32 {
	return int32(interfaces.CreateSuccess)
}

// Read reads the packet data from the given reader
func (p *CreateSuccess) Read(r interfaces.Reader) error {
	var err error
	p.ObjectId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.CharId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Stats, err = r.ReadString()
	return err
}

// Write writes the packet data to the given writer
func (p *CreateSuccess) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.ObjectId); err != nil {
		return err
	}

	if err := w.WriteInt32(p.CharId); err != nil {
		return err
	}

	return w.WriteString(p.Stats)
}

// String returns a string representation of the packet
func (p *CreateSuccess) String() string {
	return "CreateSuccess"
}

// HasNulls checks if any fields in the packet are null
func (p *CreateSuccess) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *CreateSuccess) Structure() string {
	return "CreateSuccess"
}
