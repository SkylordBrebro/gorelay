﻿package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// Hello represents a packet for client hello handshake
type Hello struct {
	*packets.BasePacket
	GameID               int32
	BuildVersion         string
	AccessToken          string
	KeyTime              int32
	Key                  []byte
	GameNet              string
	PlayPlatform         string
	PlatformToken        string
	ClientToken          string
	ClientIdentification string
}

// NewHello creates a new Hello packet
func NewHello() *Hello {
	return &Hello{
		BasePacket: packets.NewPacket(interfaces.Hello, byte(interfaces.Hello)),
	}
}

// Type returns the packet type
func (p *Hello) Type() interfaces.PacketType {
	return interfaces.Hello
}

// Read reads the packet data from a PacketReader
func (p *Hello) Read(r *packets.PacketReader) error {
	var err error
	p.GameID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.BuildVersion, err = r.ReadString()
	if err != nil {
		return err
	}
	p.AccessToken, err = r.ReadString()
	if err != nil {
		return err
	}
	p.KeyTime, err = r.ReadInt32()
	if err != nil {
		return err
	}
	keyLen, err := r.ReadInt16()
	if err != nil {
		return err
	}
	p.Key, err = r.ReadBytes(int(keyLen))
	if err != nil {
		return err
	}
	p.GameNet, err = r.ReadString()
	if err != nil {
		return err
	}
	p.PlayPlatform, err = r.ReadString()
	if err != nil {
		return err
	}
	p.PlatformToken, err = r.ReadString()
	if err != nil {
		return err
	}
	p.ClientToken, err = r.ReadString()
	if err != nil {
		return err
	}
	p.ClientIdentification, err = r.ReadString()
	return err
}

// Write writes the packet data to a PacketWriter
func (p *Hello) Write(w *packets.PacketWriter) error {
	if err := w.WriteInt32(p.GameID); err != nil {
		return err
	}
	if err := w.WriteString(p.BuildVersion); err != nil {
		return err
	}
	if err := w.WriteString(p.AccessToken); err != nil {
		return err
	}
	if err := w.WriteInt32(p.KeyTime); err != nil {
		return err
	}
	if err := w.WriteInt16(int16(len(p.Key))); err != nil {
		return err
	}
	if err := w.WriteBytes(p.Key); err != nil {
		return err
	}
	if err := w.WriteString(p.GameNet); err != nil {
		return err
	}
	if err := w.WriteString(p.PlayPlatform); err != nil {
		return err
	}
	if err := w.WriteString(p.PlatformToken); err != nil {
		return err
	}
	if err := w.WriteString(p.ClientToken); err != nil {
		return err
	}
	return w.WriteString(p.ClientIdentification)
}
