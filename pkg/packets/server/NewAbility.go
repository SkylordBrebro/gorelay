package server

import (
	"gorelay/pkg/packets/interfaces"
)

// NewAbility represents the server packet for new ability notification
type NewAbility struct {
	AbilityType int32
}

// Type returns the packet type for NewAbility
func (p *NewAbility) Type() interfaces.PacketType {
	return interfaces.NewAbility
}

// Read reads the packet data from the provided reader
func (p *NewAbility) Read(r interfaces.Reader) error {
	var err error
	p.AbilityType, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the provided writer
func (p *NewAbility) Write(w interfaces.Writer) error {
	return w.WriteInt32(p.AbilityType)
}

func (p *NewAbility) ID() int32 {
	return int32(interfaces.NewAbility)
}