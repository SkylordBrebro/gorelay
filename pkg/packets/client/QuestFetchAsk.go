package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// QuestFetchAsk represents a packet for requesting quest information
type QuestFetchAsk struct {
	*packets.BasePacket
}

// NewQuestFetchAsk creates a new QuestFetchAsk packet
func NewQuestFetchAsk() *QuestFetchAsk {
	return &QuestFetchAsk{
		BasePacket: packets.NewPacket(interfaces.QuestFetchAsk, byte(interfaces.QuestFetchAsk)),
	}
}

// Type returns the packet type
func (q *QuestFetchAsk) Type() interfaces.PacketType {
	return interfaces.QuestFetchAsk
}

// Read reads the packet data from a PacketReader
func (q *QuestFetchAsk) Read(r *packets.PacketReader) error {
	return nil
}

// Write writes the packet data to a PacketWriter
func (q *QuestFetchAsk) Write(w *packets.PacketWriter) error {
	return nil
}
