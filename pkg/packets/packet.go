package packets

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
	"reflect"
)

// Packet defines the interface that all packets must implement
type Packet interface {
	Type() interfaces.PacketType
	ID() int32
	Read(r interfaces.Reader) error
	Write(w interfaces.Writer) error
	String() string
	HasNulls() bool
	Structure() string
}

// BasePacket represents a base network packet
type BasePacket struct {
	Send     bool
	PacketID byte
	Extra    []byte
	data     []byte
}

// Type returns the packet type
func (p *BasePacket) Type() interfaces.PacketType {
	return interfaces.Unknown
}

// Read reads packet data from a packet reader
func (p *BasePacket) Read(r interfaces.Reader) error {
	data, err := r.ReadBytes(int(r.RemainingBytes()))
	if err != nil {
		return err
	}
	p.data = data
	return nil
}

// Write writes packet data to a packet writer
func (p *BasePacket) Write(w interfaces.Writer) error {
	return w.WriteBytes(p.data)
}

// NewPacket creates a new packet instance of the specified type
func NewPacket(packetType interfaces.PacketType, packetID byte) *BasePacket {
	return &BasePacket{
		Send:     true,
		PacketID: packetID,
	}
}

// NewPacketFromData creates a new packet from raw byte data
func NewPacketFromData(data []byte) (*BasePacket, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("packet data too short")
	}

	reader := NewPacketReader(data)

	// Read length
	_, err := reader.ReadInt32()
	if err != nil {
		return nil, fmt.Errorf("failed to read packet length: %v", err)
	}

	// Read ID
	id, err := reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read packet ID: %v", err)
	}

	packet := &BasePacket{
		Send:     true,
		PacketID: id,
	}

	// Read remaining data
	if err := packet.Read(reader); err != nil {
		return nil, fmt.Errorf("failed to read packet data: %v", err)
	}

	return packet, nil
}

// String returns a string representation of the packet
func (p *BasePacket) String() string {
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)

	str := fmt.Sprintf("%s(%d) Packet Instance", p.Type(), p.PacketID)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Uint8 {
			// Format byte slices as hex
			str += fmt.Sprintf("\n\t%s => %X", field.Name, value.Interface())
		} else {
			str += fmt.Sprintf("\n\t%s => %v", field.Name, value.Interface())
		}
	}

	return str
}

// HasNulls checks if any fields in the packet are null
func (p *BasePacket) HasNulls() bool {
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if value.IsNil() {
			fmt.Printf("Packet %d has null field: %s\n", p.PacketID, field.Name)
			return true
		}
	}

	return false
}

// Structure returns a string representation of the packet structure
func (p *BasePacket) Structure() string {
	t := reflect.TypeOf(p)
	str := fmt.Sprintf("%s [%d] \nPacket Structure:\n{", p.Type(), p.PacketID)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		str += fmt.Sprintf("\n  %s => %s", field.Name, field.Type.Name())
	}

	str += "\n}"
	return str
}

// ID returns the packet ID
func (p *BasePacket) ID() int32 {
	return int32(p.PacketID)
}

// EncodePacket encodes a packet into a byte array ready for transmission
// It adds the packet length and ID, then writes the packet data
func EncodePacket(packet Packet) ([]byte, error) {
	// Create a new packet writer
	writer := NewPacketWriter()

	// Get the packet ID
	packetID := packet.ID()

	// Create a temporary writer to measure the packet size
	tempWriter := NewPacketWriter()

	// Write the packet ID to the temp writer
	if err := tempWriter.WriteByte(byte(packetID)); err != nil {
		return nil, fmt.Errorf("failed to write packet ID: %v", err)
	}

	// Write the packet data to the temp writer
	if err := packet.Write(tempWriter); err != nil {
		return nil, fmt.Errorf("failed to write packet data: %v", err)
	}

	// Get the packet data
	packetData := tempWriter.Bytes()

	// Write the packet length (including ID byte) to the main writer
	if err := writer.WriteInt32(int32(len(packetData))); err != nil {
		return nil, fmt.Errorf("failed to write packet length: %v", err)
	}

	// Write the packet data (including ID) to the main writer
	if err := writer.WriteBytes(packetData); err != nil {
		return nil, fmt.Errorf("failed to write packet data: %v", err)
	}

	// Return the encoded packet
	return writer.Bytes(), nil
}
