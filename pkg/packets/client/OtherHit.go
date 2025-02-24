// Decompiled with JetBrains decompiler
// Type: prankster.Proxy.Networking.Packets.Client.OtherHit
// Assembly: prankster, Version=1.0.0.1, Culture=neutral, PublicKeyToken=null
// MVID: 674C3C29-3FFB-46FB-A4BE-03322F13731C
// Assembly location: \\hv\e$\rotmg\multisource\pranksterREAL-cleaned.exe

package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// OtherHit represents a packet for other entity hits
type OtherHit struct {
	*packets.BasePacket
	Time     int32
	BulletID uint16
	ObjectID int32
	TargetID int32
}

// NewOtherHit creates a new OtherHit packet
func NewOtherHit() *OtherHit {
	return &OtherHit{
		BasePacket: packets.NewPacket(interfaces.OtherHit, byte(interfaces.OtherHit)),
	}
}

// Type returns the packet type
func (p *OtherHit) Type() interfaces.PacketType {
	return interfaces.OtherHit
}

// Read reads the packet data from a PacketReader
func (p *OtherHit) Read(r *packets.PacketReader) error {
	var err error
	p.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.BulletID, err = r.ReadUInt16()
	if err != nil {
		return err
	}
	p.ObjectID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.TargetID, err = r.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *OtherHit) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.Time); err != nil {
		return err
	}
	if err := w.WriteUInt16(p.BulletID); err != nil {
		return err
	}
	if err := w.WriteInt32(p.ObjectID); err != nil {
		return err
	}
	return w.WriteInt32(p.TargetID)
}
