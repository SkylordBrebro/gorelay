package server

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
	"strconv"
	"strings"
)

// Reconnect represents the server packet for reconnection information
type Reconnect struct {
	Name      string
	Host      string
	Port      uint16
	GameId    int32
	KeyTime   int32
	Key       []byte
	AliveTime int // special use, don't use for packet
}

// Type returns the packet type for Reconnect
func (p *Reconnect) Type() interfaces.PacketType {
	return interfaces.Reconnect
}

// Read reads the packet data from the provided reader
func (p *Reconnect) Read(r interfaces.Reader) error {
	var err error

	// Read Name
	p.Name, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read Host
	p.Host, err = r.ReadString()
	if err != nil {
		return err
	}

	// Read Port
	p.Port, err = r.ReadUInt16()
	if err != nil {
		return err
	}

	// Read GameId
	p.GameId, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read KeyTime
	p.KeyTime, err = r.ReadInt32()
	if err != nil {
		return err
	}

	// Read Key length
	keyLength, err := r.ReadInt16()
	if err != nil {
		return err
	}

	// Read Key
	p.Key, err = r.ReadBytes(int(keyLength))
	if err != nil {
		return err
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *Reconnect) Write(w interfaces.Writer) error {
	var err error

	// Write Name
	err = w.WriteString(p.Name)
	if err != nil {
		return err
	}

	// Write Host
	err = w.WriteString(p.Host)
	if err != nil {
		return err
	}

	// Write Port
	err = w.WriteUInt16(p.Port)
	if err != nil {
		return err
	}

	// Write GameId
	err = w.WriteInt32(p.GameId)
	if err != nil {
		return err
	}

	// Write KeyTime
	err = w.WriteInt32(p.KeyTime)
	if err != nil {
		return err
	}

	// Write Key length
	err = w.WriteInt16(int16(len(p.Key)))
	if err != nil {
		return err
	}

	// Write Key
	err = w.WriteBytes(p.Key)
	if err != nil {
		return err
	}

	return nil
}

// Serialize converts the Reconnect packet to a string representation
func (p *Reconnect) Serialize() string {
	return fmt.Sprintf("%s|%s|%d|%d|%d|%s",
		p.Name, p.Host, p.Port, p.GameId, p.KeyTime, ByteArrayToHexString(p.Key))
}

// Deserialize creates a Reconnect packet from a string representation
func Deserialize(input string) *Reconnect {
	parts := strings.Split(input, "|")
	reconnect := &Reconnect{
		Name: parts[0],
		Host: parts[1],
	}

	if port, err := strconv.ParseUint(parts[2], 10, 16); err == nil {
		reconnect.Port = uint16(port)
	}

	if gameId, err := strconv.ParseInt(parts[3], 10, 32); err == nil {
		reconnect.GameId = int32(gameId)
	}

	if keyTime, err := strconv.ParseInt(parts[4], 10, 32); err == nil {
		reconnect.KeyTime = int32(keyTime)
	}

	if len(parts) == 6 {
		reconnect.Key = HexStringToByteArray(parts[5])
	}

	return reconnect
}

// ByteArrayToHexString converts a byte array to a hex string
func ByteArrayToHexString(data []byte) string {
	if data == nil {
		return ""
	}

	var sb strings.Builder
	for _, b := range data {
		sb.WriteString(fmt.Sprintf("%02x", b))
	}
	return sb.String()
}

// HexStringToByteArray converts a hex string to a byte array
func HexStringToByteArray(hex string) []byte {
	if hex == "" {
		return nil
	}

	if len(hex)%2 != 0 {
		hex = "0" + hex
	}

	result := make([]byte, len(hex)/2)
	for i := 0; i < len(hex); i += 2 {
		if val, err := strconv.ParseUint(hex[i:i+2], 16, 8); err == nil {
			result[i/2] = byte(val)
		}
	}

	return result
}
