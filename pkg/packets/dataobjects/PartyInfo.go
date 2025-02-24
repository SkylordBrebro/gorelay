package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// PartyInfo represents information about a party
type PartyInfo struct {
	ID       int32
	Players  []*PartyPlayer
	Leader   int32
	Activity PartyActivity
	Privacy  PartyPrivacy
}

// NewPartyInfo creates a new PartyInfo instance
func NewPartyInfo() *PartyInfo {
	return &PartyInfo{
		Players: make([]*PartyPlayer, 0),
	}
}

// Read reads the party info from a Reader
func (p *PartyInfo) Read(r interfaces.Reader) error {
	var err error
	p.ID, err = r.ReadInt32()
	if err != nil {
		return err
	}

	playerCount, err := r.ReadInt16()
	if err != nil {
		return err
	}

	p.Players = make([]*PartyPlayer, playerCount)
	for i := 0; i < int(playerCount); i++ {
		player := NewPartyPlayer()
		if err := player.Read(r); err != nil {
			return err
		}
		p.Players[i] = player
	}

	p.Leader, err = r.ReadInt32()
	if err != nil {
		return err
	}

	activityByte, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Activity = PartyActivity(activityByte)

	privacyByte, err := r.ReadByte()
	if err != nil {
		return err
	}
	p.Privacy = PartyPrivacy(privacyByte)

	return nil
}

// Write writes the party info to a Writer
func (p *PartyInfo) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(p.ID); err != nil {
		return err
	}

	if err := w.WriteInt16(int16(len(p.Players))); err != nil {
		return err
	}

	for _, player := range p.Players {
		if err := player.Write(w); err != nil {
			return err
		}
	}

	if err := w.WriteInt32(p.Leader); err != nil {
		return err
	}

	if err := w.WriteByte(byte(p.Activity)); err != nil {
		return err
	}

	return w.WriteByte(byte(p.Privacy))
}

// Clone creates a copy of the PartyInfo
func (p *PartyInfo) Clone() DataObject {
	players := make([]*PartyPlayer, len(p.Players))
	for i, player := range p.Players {
		players[i] = player.Clone().(*PartyPlayer)
	}

	return &PartyInfo{
		ID:       p.ID,
		Players:  players,
		Leader:   p.Leader,
		Activity: p.Activity,
		Privacy:  p.Privacy,
	}
}

// String returns a string representation of the PartyInfo
func (p *PartyInfo) String() string {
	return fmt.Sprintf("{ Id=%d, Players=%d, Leader=%d, Activity=%v, Privacy=%v }",
		p.ID, len(p.Players), p.Leader, p.Activity, p.Privacy)
}
