package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// GoToQuestRoom represents a packet for going to quest room
type GoToQuestRoom struct {
	*packets.BasePacket
}

// NewGoToQuestRoom creates a new GoToQuestRoom packet
func NewGoToQuestRoom() *GoToQuestRoom {
	return &GoToQuestRoom{
		BasePacket: packets.NewPacket(interfaces.GoToQuestRoom, byte(interfaces.GoToQuestRoom)),
	}
}

// Type returns the packet type
func (p *GoToQuestRoom) Type() interfaces.PacketType {
	return interfaces.GoToQuestRoom
}

// Read reads the packet data from a PacketReader
func (p *GoToQuestRoom) Read(r *packets.PacketReader) error {
	return nil
}

// Write writes the packet data to a PacketWriter
func (p *GoToQuestRoom) Write(w *packets.PacketWriter) error {
	return nil
}
