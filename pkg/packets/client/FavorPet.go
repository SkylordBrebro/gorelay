// Decompiled with JetBrains decompiler
// Type: prankster.Proxy.Networking.Packets.Client.FavorPet
// Assembly: prankster, Version=1.0.0.1, Culture=neutral, PublicKeyToken=null
// MVID: 674C3C29-3FFB-46FB-A4BE-03322F13731C
// Assembly location: \\hv\e$\rotmg\multisource\pranksterREAL-cleaned.exe

package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// FavorPet represents a packet for favoriting a pet
type FavorPet struct {
	*packets.BasePacket
	PetID int32
}

// NewFavorPet creates a new FavorPet packet
func NewFavorPet() *FavorPet {
	return &FavorPet{
		BasePacket: packets.NewPacket(interfaces.FavorPet, byte(interfaces.FavorPet)),
	}
}

// Type returns the packet type
func (p *FavorPet) Type() interfaces.PacketType {
	return interfaces.FavorPet
}

// Read reads the packet data from a PacketReader
func (p *FavorPet) Read(r *packets.PacketReader) error {
	var err error
	p.PetID, err = r.ReadInt32()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *FavorPet) Write(w *packets.PacketWriter) error {
	return w.WriteInt32(p.PetID)
}
