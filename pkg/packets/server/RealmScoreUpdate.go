package server

import (
	"gorelay/pkg/packets/interfaces"
)

// RealmScoreUpdate represents the server packet for realm score updates
type RealmScoreUpdate struct {
	Score int32 // min 0 max 300000
}

// Type returns the packet type for RealmScoreUpdate
func (p *RealmScoreUpdate) Type() interfaces.PacketType {
	return interfaces.RealmScoreUpdate
}

// Read reads the packet data from the provided reader
func (p *RealmScoreUpdate) Read(r interfaces.Reader) error {
	var err error
	p.Score, err = r.ReadInt32()
	return err
}

// Write writes the packet data to the provided writer
func (p *RealmScoreUpdate) Write(w interfaces.Writer) error {
	return w.WriteInt32(p.Score)
}

func (p *RealmScoreUpdate) ID() int32 {
	return int32(interfaces.RealmScoreUpdate)
}