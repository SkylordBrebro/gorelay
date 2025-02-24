package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// PartyPlayer represents a player in a party
type PartyPlayer struct {
	Name     string
	ObjectID int32
	Level    int16
	Class    int16
}

// NewPartyPlayer creates a new PartyPlayer instance
func NewPartyPlayer() *PartyPlayer {
	return &PartyPlayer{}
}

// Read reads the party player data from a Reader
func (p *PartyPlayer) Read(r interfaces.Reader) error {
	var err error
	p.Name, err = r.ReadString()
	if err != nil {
		return err
	}
	p.ObjectID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	p.Level, err = r.ReadInt16()
	if err != nil {
		return err
	}
	p.Class, err = r.ReadInt16()
	return err
}

// Write writes the party player data to a PacketWriter
func (p *PartyPlayer) Write(w interfaces.Writer) error {
	if err := w.WriteString(p.Name); err != nil {
		return err
	}
	if err := w.WriteInt32(p.ObjectID); err != nil {
		return err
	}
	if err := w.WriteInt16(p.Level); err != nil {
		return err
	}
	return w.WriteInt16(p.Class)
}

// Clone creates a copy of the PartyPlayer
func (p *PartyPlayer) Clone() DataObject {
	return &PartyPlayer{
		Name:     p.Name,
		ObjectID: p.ObjectID,
		Level:    p.Level,
		Class:    p.Class,
	}
}

// String returns a string representation of the PartyPlayer
func (p *PartyPlayer) String() string {
	return fmt.Sprintf("{ Name=%s, ObjectId=%d, Level=%d, Class=%d }",
		p.Name, p.ObjectID, p.Level, p.Class)
}
