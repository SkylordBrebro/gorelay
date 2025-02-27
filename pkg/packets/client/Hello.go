package client

import (
	"encoding/hex"
	"fmt"
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

// Read reads the packet data from a Reader
func (p *Hello) Read(r interfaces.Reader) error {
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

// Write writes the packet data to a Writer
func (p *Hello) Write(w interfaces.Writer) error {
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

// String returns a string representation of the Hello struct
func (p *Hello) ToString() string {
	return fmt.Sprintf("GameID: %d, BuildVersion: %s, AccessToken: %s, KeyTime: %d, Key: %s, GameNet: %s, PlayPlatform: %s, PlatformToken: %s, ClientToken: %s, ClientIdentification: %s",
		p.GameID, p.BuildVersion, p.AccessToken, p.KeyTime, hex.EncodeToString(p.Key), p.GameNet, p.PlayPlatform, p.PlatformToken, p.ClientToken, p.ClientIdentification)
}

// String returns a string representation of the packet
func (p *Hello) String() string {
	return p.ToString()
}

// HasNulls checks if any fields in the packet are null
func (p *Hello) HasNulls() bool {
	return false
}

// Structure returns a string representation of the packet structure
func (p *Hello) Structure() string {
	return fmt.Sprintf("Hello Packet (ID=%d)", p.ID())
}

// ID returns the packet ID
func (p *Hello) ID() int32 {
	return int32(interfaces.Hello)
}

