package server

import (
	"gorelay/pkg/packets/interfaces"
)

// UnlockNewSlot represents the server packet for unlocking a new character slot
type UnlockNewSlot struct {
	UnlockType int32
}

// Type returns the packet type for UnlockNewSlot
func (p *UnlockNewSlot) Type() interfaces.PacketType {
	return interfaces.UnlockNewSlot
}

// Read reads the packet data from the provided reader
func (p *UnlockNewSlot) Read(r interfaces.Reader) error {
	var err error
	p.UnlockType, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the provided writer
func (p *UnlockNewSlot) Write(w interfaces.Writer) error {
	return w.WriteInt32(p.UnlockType)
}
