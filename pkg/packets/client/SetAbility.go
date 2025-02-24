package client

import (
	"gorelay/pkg/packets/interfaces"
)

// SetAbility represents a client-side ability setting packet
type SetAbility struct {
	AbilityID int32
	Status    int8
}

// Type returns the packet type for SetAbility
func (p *SetAbility) Type() interfaces.PacketType {
	return interfaces.SetAbility
}

// Read reads the packet data from the given reader
func (p *SetAbility) Read(r interfaces.Reader) error {
	var err error
	p.AbilityID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	status, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Status = int8(status)
	return nil
}

// Write writes the packet data to the given writer
func (p *SetAbility) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.AbilityID); err != nil {
		return err
	}
	return w.WriteByte(byte(p.Status))
}
