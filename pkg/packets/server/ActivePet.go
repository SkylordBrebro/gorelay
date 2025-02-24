package server

import (
	"gorelay/pkg/packets/interfaces"
)

// ActivePet represents a server-side active pet packet
type ActivePet struct {
	PetId int32
}

// Type returns the packet type for ActivePet
func (p *ActivePet) Type() interfaces.PacketType {
	return interfaces.ActivePet
}

// ID returns the packet ID
func (p *ActivePet) ID() int32 {
	return int32(interfaces.ActivePet)
}

// Read reads the packet data from the given reader
func (p *ActivePet) Read(r interfaces.Reader) error {
	var err error
	p.PetId, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the given writer
func (p *ActivePet) Write(w interfaces.Writer) error {
	return w.WriteInt32(p.PetId)
}

// String returns a string representation of the packet
func (p *ActivePet) String() string {
	return "ActivePet"
}

// HasNulls checks if any fields in the packet are null
func (p *ActivePet) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *ActivePet) Structure() string {
	return "ActivePet"
}
