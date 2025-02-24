package packets

import (
	"fmt"
	"reflect"
)

// PacketHandler manages packet routing and handling
type PacketHandler struct {
	handlers map[int]PacketHandlerFunc
	packets  map[int]reflect.Type
}

type PacketHandlerFunc func(data []byte) error

// NewPacketHandler creates a new packet handler instance
func NewPacketHandler() *PacketHandler {
	return &PacketHandler{
		handlers: make(map[int]PacketHandlerFunc),
		packets:  make(map[int]reflect.Type),
	}
}

// RegisterHandler registers a handler for a specific packet type
func (ph *PacketHandler) RegisterHandler(packetID int, handler PacketHandlerFunc) {
	ph.handlers[packetID] = handler
}

// RegisterPacket registers a packet type with its ID
func (ph *PacketHandler) RegisterPacket(packetID int, packet interface{}) error {
	t := reflect.TypeOf(packet)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	ph.packets[packetID] = t
	return nil
}

// HandlePacket processes an incoming packet
func (ph *PacketHandler) HandlePacket(packetID int, data []byte) error {
	if handler, ok := ph.handlers[packetID]; ok {
		return handler(data)
	}
	return fmt.Errorf("no handler registered for packet ID %d", packetID)
}

// GetPacketType returns the type for a given packet ID
func (ph *PacketHandler) GetPacketType(packetID int) (reflect.Type, error) {
	if t, ok := ph.packets[packetID]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("no packet type registered for ID %d", packetID)
}

// ClearHandlers removes all registered handlers
func (ph *PacketHandler) ClearHandlers() {
	ph.handlers = make(map[int]PacketHandlerFunc)
}

// ClearPackets removes all registered packet types
func (ph *PacketHandler) ClearPackets() {
	ph.packets = make(map[int]reflect.Type)
}
