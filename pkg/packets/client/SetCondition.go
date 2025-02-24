package client

import (
	"gorelay/pkg/packets/interfaces"
)

// SetCondition represents a client-side condition setting packet
type SetCondition struct {
	ConditionEffect   byte
	ConditionDuration float32
}

// Type returns the packet type for SetCondition
func (p *SetCondition) Type() interfaces.PacketType {
	return interfaces.SetCondition
}

// Read reads the packet data from the given reader
func (p *SetCondition) Read(r interfaces.Reader) error {
	var err error
	p.ConditionEffect, err = r.ReadByte()
	if err != nil {
		return err
	}
	p.ConditionDuration, err = r.ReadFloat32()
	return err
}

// Write writes the packet data to the given writer
func (p *SetCondition) Write(w interfaces.Writer) error {
	if err := w.WriteByte(p.ConditionEffect); err != nil {
		return err
	}
	return w.WriteFloat32(p.ConditionDuration)
}
