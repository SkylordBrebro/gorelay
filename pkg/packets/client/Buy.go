// Decompiled with JetBrains decompiler
// Type: prankster.Proxy.Networking.Packets.Client.Buy
// Assembly: prankster, Version=1.0.0.1, Culture=neutral, PublicKeyToken=null
// MVID: 674C3C29-3FFB-46FB-A4BE-03322F13731C
// Assembly location: \\hv\e$\rotmg\multisource\pranksterREAL-cleaned.exe

package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// Buy represents a purchase packet
type Buy struct {
	*packets.BasePacket
	ObjectID int32
	Quantity int32
}

// NewBuy creates a new Buy packet
func NewBuy() *Buy {
	return &Buy{
		BasePacket: packets.NewPacket(interfaces.Buy, byte(interfaces.Buy)),
	}
}

// Type returns the packet type
func (p *Buy) Type() interfaces.PacketType {
	return interfaces.Buy
}

// Read reads the packet data from the reader
func (p *Buy) Read(r *packets.PacketReader) error {
	var err error
	p.ObjectID, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.Quantity, err = r.ReadInt32()
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the writer
func (p *Buy) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.ObjectID); err != nil {
		return err
	}

	return w.WriteInt32(p.Quantity)
}
