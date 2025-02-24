package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// QueueCancel represents a packet for canceling a queue
type QueueCancel struct {
	*packets.BasePacket
}

// NewQueueCancel creates a new QueueCancel packet
func NewQueueCancel() *QueueCancel {
	return &QueueCancel{
		BasePacket: packets.NewPacket(interfaces.QueueCancel, byte(interfaces.QueueCancel)),
	}
}

// Type returns the packet type
func (q *QueueCancel) Type() interfaces.PacketType {
	return interfaces.QueueCancel
}

// Read reads the packet data from a PacketReader
func (q *QueueCancel) Read(r *packets.PacketReader) error {
	return nil
}

// Write writes the packet data to a PacketWriter
func (q *QueueCancel) Write(w *packets.PacketWriter) error {
	return nil
}
