package client

import (
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/interfaces"
)

// RequestTrade represents a packet for requesting a trade
type RequestTrade struct {
	*packets.BasePacket
	Name string
}

// NewRequestTrade creates a new RequestTrade packet
func NewRequestTrade() *RequestTrade {
	return &RequestTrade{
		BasePacket: packets.NewPacket(interfaces.RequestTrade, byte(interfaces.RequestTrade)),
	}
}

// Type returns the packet type
func (r *RequestTrade) Type() interfaces.PacketType {
	return interfaces.RequestTrade
}

// Read reads the packet data from a PacketReader
func (r *RequestTrade) Read(reader *packets.PacketReader) error {
	var err error
	r.Name, err = reader.ReadString()
	return err
}

// Write writes the packet data to a PacketWriter
func (r *RequestTrade) Write(writer *packets.PacketWriter) error {
	return writer.WriteString(r.Name)
}
