package packets

import (
	"bytes"
	"encoding/binary"
	"math"
)

// PacketWriter handles writing binary data to packets
type PacketWriter struct {
	buffer *bytes.Buffer
}

// NewPacketWriter creates a new packet writer
func NewPacketWriter() *PacketWriter {
	return &PacketWriter{
		buffer: bytes.NewBuffer(nil),
	}
}

// WriteInt16 writes a network-ordered int16
func (pw *PacketWriter) WriteInt16(value int16) error {
	return binary.Write(pw.buffer, binary.BigEndian, value)
}

// WriteUInt16 writes a network-ordered uint16
func (pw *PacketWriter) WriteUInt16(value uint16) error {
	return binary.Write(pw.buffer, binary.BigEndian, value)
}

// WriteUInt32 writes a network-ordered uint32
func (pw *PacketWriter) WriteUInt32(value uint32) error {
	return binary.Write(pw.buffer, binary.BigEndian, value)
}

// WriteInt32 writes a network-ordered int32
func (pw *PacketWriter) WriteInt32(value int32) error {
	return binary.Write(pw.buffer, binary.BigEndian, value)
}

// WriteFloat32 writes a network-ordered float32
func (pw *PacketWriter) WriteFloat32(value float32) error {
	return binary.Write(pw.buffer, binary.BigEndian, math.Float32bits(value))
}

// WriteString writes a length-prefixed UTF-8 string
func (pw *PacketWriter) WriteString(value string) error {
	data := []byte(value)
	if err := pw.WriteInt16(int16(len(data))); err != nil {
		return err
	}
	_, err := pw.buffer.Write(data)
	return err
}

// WriteUTF32String writes a 32-bit length-prefixed UTF-8 string
func (pw *PacketWriter) WriteUTF32String(value string) error {
	data := []byte(value)
	if err := pw.WriteInt32(int32(len(data))); err != nil {
		return err
	}
	_, err := pw.buffer.Write(data)
	return err
}

// WriteCompressedInt writes a compressed integer value
func (pw *PacketWriter) WriteCompressedInt(value int) error {
	isNegative := value < 0
	if isNegative {
		value = -value
	}

	num := uint(value)
	firstByte := byte(num & 63)
	if isNegative {
		firstByte |= 64
	}
	num >>= 6

	if num > 0 {
		firstByte |= 128
	}

	if err := pw.buffer.WriteByte(firstByte); err != nil {
		return err
	}

	for num > 0 {
		b := byte(num & 127)
		num >>= 7
		if num > 0 {
			b |= 128
		}
		if err := pw.buffer.WriteByte(b); err != nil {
			return err
		}
	}

	return nil
}

// WriteByte writes a single byte
func (pw *PacketWriter) WriteByte(value byte) error {
	return pw.buffer.WriteByte(value)
}

// WriteBytes writes a byte slice
func (pw *PacketWriter) WriteBytes(data []byte) error {
	_, err := pw.buffer.Write(data)
	return err
}

// Bytes returns the written bytes
func (pw *PacketWriter) Bytes() []byte {
	return pw.buffer.Bytes()
}

// BlockCopyInt32 copies a network-ordered int32 into a byte slice at index 0
func BlockCopyInt32(data []byte, value int32) {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(value))
	data[0] = bytes[0]
	data[1] = bytes[1]
	data[2] = bytes[2]
	data[3] = bytes[3]
}

// WriteBool writes a boolean value
func (pw *PacketWriter) WriteBool(value bool) error {
	if value {
		return pw.WriteByte(1)
	}
	return pw.WriteByte(0)
}
