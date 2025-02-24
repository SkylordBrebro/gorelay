package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// PacketReader provides utilities for reading packet data
type PacketReader struct {
	reader io.Reader
}

// NewPacketReader creates a new packet reader
func NewPacketReader(data []byte) *PacketReader {
	return &PacketReader{
		reader: bytes.NewReader(data),
	}
}

// ReadByte reads a single byte
func (r *PacketReader) ReadByte() (byte, error) {
	var b byte
	err := binary.Read(r.reader, binary.BigEndian, &b)
	return b, err
}

// ReadInt32 reads a 32-bit integer
func (r *PacketReader) ReadInt32() (int32, error) {
	var n int32
	err := binary.Read(r.reader, binary.BigEndian, &n)
	return n, err
}

// ReadUInt32 reads a 32-bit unsigned integer
func (r *PacketReader) ReadUInt32() (uint32, error) {
	var n uint32
	err := binary.Read(r.reader, binary.BigEndian, &n)
	return n, err
}

// ReadFloat32 reads a 32-bit float
func (r *PacketReader) ReadFloat32() (float32, error) {
	var f float32
	err := binary.Read(r.reader, binary.BigEndian, &f)
	return f, err
}

// ReadString reads a string
func (r *PacketReader) ReadString() (string, error) {
	return ReadUTF(r.reader)
}

// PacketWriter provides utilities for writing packet data
type PacketWriter struct {
	buffer *bytes.Buffer
}

// NewPacketWriter creates a new packet writer
func NewPacketWriter() *PacketWriter {
	return &PacketWriter{
		buffer: bytes.NewBuffer(nil),
	}
}

// WriteByte writes a single byte
func (w *PacketWriter) WriteByte(b byte) error {
	return binary.Write(w.buffer, binary.BigEndian, b)
}

// WriteInt32 writes a 32-bit integer
func (w *PacketWriter) WriteInt32(n int32) error {
	return binary.Write(w.buffer, binary.BigEndian, n)
}

// WriteUInt32 writes a 32-bit unsigned integer
func (w *PacketWriter) WriteUInt32(n uint32) error {
	return binary.Write(w.buffer, binary.BigEndian, n)
}

// WriteFloat32 writes a 32-bit float
func (w *PacketWriter) WriteFloat32(f float32) error {
	return binary.Write(w.buffer, binary.BigEndian, f)
}

// WriteString writes a string
func (w *PacketWriter) WriteString(s string) error {
	return WriteUTF(w.buffer, s)
}

// Bytes returns the written bytes
func (w *PacketWriter) Bytes() []byte {
	return w.buffer.Bytes()
}

// CompressInt32 compresses a 32-bit integer into a variable-length format
func CompressInt32(n int32) []byte {
	var buf bytes.Buffer
	for n >= 0x80 {
		buf.WriteByte(byte(n) | 0x80)
		n >>= 7
	}
	buf.WriteByte(byte(n))
	return buf.Bytes()
}

// DecompressInt32 decompresses a variable-length integer
func DecompressInt32(r io.Reader) (int32, error) {
	var result int32
	var shift uint
	for {
		var b byte
		if err := binary.Read(r, binary.BigEndian, &b); err != nil {
			return 0, fmt.Errorf("failed to read byte: %v", err)
		}
		result |= int32(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		if shift >= 32 {
			return 0, fmt.Errorf("integer too large")
		}
	}
	return result, nil
}

// PacketDataReader provides utilities for reading game-specific data types
type PacketDataReader struct {
	*PacketReader
}

// NewPacketDataReader creates a new packet data reader
func NewPacketDataReader(data []byte) *PacketDataReader {
	return &PacketDataReader{
		PacketReader: NewPacketReader(data),
	}
}

// ReadWorldPos reads a world position
func (r *PacketDataReader) ReadWorldPos() (float32, float32, error) {
	x, err := r.ReadFloat32()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read X coordinate: %v", err)
	}
	y, err := r.ReadFloat32()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read Y coordinate: %v", err)
	}
	return x, y, nil
}

// ReadAngle reads an angle in radians
func (r *PacketDataReader) ReadAngle() (float32, error) {
	raw, err := r.ReadFloat32()
	if err != nil {
		return 0, fmt.Errorf("failed to read angle: %v", err)
	}
	return raw * float32(math.Pi/180.0), nil
}

// PacketDataWriter provides utilities for writing game-specific data types
type PacketDataWriter struct {
	*PacketWriter
}

// NewPacketDataWriter creates a new packet data writer
func NewPacketDataWriter() *PacketDataWriter {
	return &PacketDataWriter{
		PacketWriter: NewPacketWriter(),
	}
}

// WriteWorldPos writes a world position
func (w *PacketDataWriter) WriteWorldPos(x, y float32) error {
	if err := w.WriteFloat32(x); err != nil {
		return fmt.Errorf("failed to write X coordinate: %v", err)
	}
	if err := w.WriteFloat32(y); err != nil {
		return fmt.Errorf("failed to write Y coordinate: %v", err)
	}
	return nil
}

// WriteAngle writes an angle in radians
func (w *PacketDataWriter) WriteAngle(angle float32) error {
	return w.WriteFloat32(angle * float32(180.0/math.Pi))
}
