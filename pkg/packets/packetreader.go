package packets

import (
	"bytes"
	"encoding/binary"
	"gorelay/pkg/packets/dataobjects"
	"io"
	"math"
)

// PacketReader handles reading binary data from packets
type PacketReader struct {
	reader *bytes.Reader
}

// NewPacketReader creates a new packet reader from a byte slice
func NewPacketReader(data []byte) *PacketReader {
	return &PacketReader{
		reader: bytes.NewReader(data),
	}
}

// ReadInt16 reads a network-ordered int16
func (pr *PacketReader) ReadInt16() (int16, error) {
	var value int16
	err := binary.Read(pr.reader, binary.BigEndian, &value)
	return value, err
}

// ReadUInt16 reads a network-ordered uint16
func (pr *PacketReader) ReadUInt16() (uint16, error) {
	var value uint16
	err := binary.Read(pr.reader, binary.BigEndian, &value)
	return value, err
}

// ReadInt32 reads a network-ordered int32
func (pr *PacketReader) ReadInt32() (int32, error) {
	var value int32
	err := binary.Read(pr.reader, binary.BigEndian, &value)
	return value, err
}

// ReadFloat32 reads a network-ordered float32
func (pr *PacketReader) ReadFloat32() (float32, error) {
	var value uint32
	err := binary.Read(pr.reader, binary.BigEndian, &value)
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(value), nil
}

// ReadString reads a length-prefixed UTF-8 string
func (pr *PacketReader) ReadString() (string, error) {
	length, err := pr.ReadInt16()
	if err != nil {
		return "", err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(pr.reader, data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ReadUTF32String reads a 32-bit length-prefixed UTF-8 string
func (pr *PacketReader) ReadUTF32String() (string, error) {
	length, err := pr.ReadInt32()
	if err != nil {
		return "", err
	}

	data := make([]byte, length)
	_, err = io.ReadFull(pr.reader, data)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ReadCompressedInt reads a compressed integer value
func (pr *PacketReader) ReadCompressedInt() (int, error) {
	b, err := pr.reader.ReadByte()
	if err != nil {
		return 0, err
	}

	isNegative := (b & 64) > 0
	num := int(b & 63)
	shift := uint(6)

	for (b & 128) != 0 {
		b, err = pr.reader.ReadByte()
		if err != nil {
			return 0, err
		}

		num |= int(b&127) << shift
		shift += 7
	}

	if isNegative {
		num = -num
	}

	return num, nil
}

// ReadByte reads a single byte
func (pr *PacketReader) ReadByte() (byte, error) {
	return pr.reader.ReadByte()
}

// ReadBytes reads n bytes
func (pr *PacketReader) ReadBytes(n int) ([]byte, error) {
	data := make([]byte, n)
	_, err := io.ReadFull(pr.reader, data)
	return data, err
}

// RemainingBytes returns the number of unread bytes
func (pr *PacketReader) RemainingBytes() byte {
	return byte(pr.reader.Len())
}

// ReadUInt32 reads a network-ordered uint32
func (pr *PacketReader) ReadUInt32() (uint32, error) {
	var value uint32
	err := binary.Read(pr.reader, binary.BigEndian, &value)
	return value, err
}

// ReadBool reads a boolean value
func (pr *PacketReader) ReadBool() (bool, error) {
	b, err := pr.ReadByte()
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

// ReadDataObject reads a DataObject from the packet
func (pr *PacketReader) ReadDataObject(obj dataobjects.DataObject) error {
	return obj.Read(pr)
}
