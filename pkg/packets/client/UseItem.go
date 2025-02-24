package client

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// UseItem represents a client-side use item packet
type UseItem struct {
	Time       int32
	SlotObject *dataobjects.SlotObject
	ItemUsePos *dataobjects.Location
	UseType    byte
}

// Type returns the packet type for UseItem
func (p *UseItem) Type() interfaces.PacketType {
	return interfaces.UseItem
}

// Read reads the packet data from the given reader
func (p *UseItem) Read(r interfaces.Reader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.SlotObject = dataobjects.NewSlotObject()
	if err := p.SlotObject.Read(r); err != nil {
		return err
	}

	p.ItemUsePos = dataobjects.NewLocation()
	if err := p.ItemUsePos.Read(r); err != nil {
		return err
	}

	p.UseType, err = r.ReadByte()
	return err
}

// Write writes the packet data to the given writer
func (p *UseItem) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	if err := p.SlotObject.Write(w); err != nil {
		return err
	}
	if err := p.ItemUsePos.Write(w); err != nil {
		return err
	}
	return w.WriteByte(p.UseType)
}
