package server

import (
	"gorelay/pkg/packets/interfaces"
)

// PlayersList represents the server packet for players list
type PlayersList struct {
	Players string
}

// Type returns the packet type for PlayersList
func (p *PlayersList) Type() interfaces.PacketType {
	return interfaces.PlayersList
}

// Read reads the packet data from the provided reader
func (p *PlayersList) Read(r interfaces.Reader) error {
	var err error
	p.Players, err = r.ReadString()
	return err
}

// Write writes the packet data to the provided writer
func (p *PlayersList) Write(w interfaces.Writer) error {
	return w.WriteString(p.Players)
}
