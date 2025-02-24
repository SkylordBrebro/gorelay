package server

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// AOE represents a server-side area of effect packet
type AOE struct {
	Location       *dataobjects.Location
	Radius         float32
	Damage         uint16
	Effects        byte
	EffectDuration float32
	OriginType     uint16
	Color          int32
	ArmorPierce    bool
}

// Type returns the packet type for AOE
func (p *AOE) Type() interfaces.PacketType {
	return interfaces.AOE
}

// ID returns the packet ID
func (p *AOE) ID() int32 {
	return int32(interfaces.AOE)
}

// Read reads the packet data from the given reader
func (p *AOE) Read(r interfaces.Reader) error {
	var err error

	p.Location = dataobjects.NewLocation()
	if err := p.Location.Read(r); err != nil {
		return err
	}

	p.Radius, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.Damage, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	p.Effects, err = r.ReadByte()
	if err != nil {
		return err
	}

	p.EffectDuration, err = r.ReadFloat32()
	if err != nil {
		return err
	}

	p.OriginType, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	p.Color, err = r.ReadInt32()
	if err != nil {
		return err
	}

	p.ArmorPierce, err = r.ReadBool()
	return err
}

// Write writes the packet data to the given writer
func (p *AOE) Write(w interfaces.Writer) error {
	if err := p.Location.Write(w); err != nil {
		return err
	}

	if err := w.WriteFloat32(p.Radius); err != nil {
		return err
	}

	if err := w.WriteUInt16(p.Damage); err != nil {
		return err
	}

	if err := w.WriteByte(p.Effects); err != nil {
		return err
	}

	if err := w.WriteFloat32(p.EffectDuration); err != nil {
		return err
	}

	if err := w.WriteUInt16(p.OriginType); err != nil {
		return err
	}

	if err := w.WriteInt32(p.Color); err != nil {
		return err
	}

	return w.WriteBool(p.ArmorPierce)
}

// String returns a string representation of the packet
func (p *AOE) String() string {
	return "AOE"
}

// HasNulls checks if any fields in the packet are null
func (p *AOE) HasNulls() bool {
	return p.Location == nil
}

// Structure returns a string representation of the packet structure
func (p *AOE) Structure() string {
	return "AOE"
}
