package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// EditAccountList represents a packet for editing an account list
type EditAccountList struct {
	*packets.BasePacket
	AccountListID int32
	Add           bool
	ObjectID      int32
}

// NewEditAccountList creates a new EditAccountList packet
func NewEditAccountList() *EditAccountList {
	return &EditAccountList{
		BasePacket: packets.NewPacket(interfaces.EditAccountList, byte(interfaces.EditAccountList)),
	}
}

// Type returns the packet type
func (p *EditAccountList) Type() interfaces.PacketType {
	return interfaces.EditAccountList
}

// Read reads the packet data from a PacketReader
func (p *EditAccountList) Read(r *packets.PacketReader) error {
	var err error
	p.AccountListID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Add, err = r.ReadBool()
	if err != nil {
		return err
	}
	p.ObjectID, err = r.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *EditAccountList) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.AccountListID); err != nil {
		return err
	}
	if err := w.WriteBool(p.Add); err != nil {
		return err
	}
	return w.WriteInt32(p.ObjectID)
}
