package server

import (
	"bytes"
	"gorelay/pkg/packets"
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// EffectType represents the type of effect to show
type EffectType byte

// Effect type constants
const (
	EffectBitColor    = 1
	EffectBitPos1X    = 2
	EffectBitPos1Y    = 4
	EffectBitPos2X    = 8
	EffectBitPos2Y    = 16
	EffectBitPos1     = 6  // EffectBitPos1X | EffectBitPos1Y
	EffectBitPos2     = 24 // EffectBitPos2X | EffectBitPos2Y
	EffectBitDuration = 32
	EffectBitId       = 64
	UnknownBitId      = 128
)

// ShowEffect represents the server packet for showing effects
type ShowEffect struct {
	EffectValue  byte
	EffectType   EffectType
	TargetId     int
	PosA         *dataobjects.Location
	PosB         *dataobjects.Location
	Color        *dataobjects.ARGB
	Duration     float32
	UnknownValue byte
}

// Type returns the packet type for ShowEffect
func (p *ShowEffect) Type() interfaces.PacketType {
	return interfaces.ShowEffect
}

// Read reads the packet data from the provided reader
func (p *ShowEffect) Read(r interfaces.Reader) error {
	var err error

	// Initialize default values
	p.PosA = dataobjects.NewLocation()
	p.PosB = dataobjects.NewLocation()
	p.Color = dataobjects.EmptyARGB()
	p.Duration = 1.0
	p.UnknownValue = 100

	// Read EffectValue
	p.EffectValue, err = r.ReadByte()
	if err != nil {
		return err
	}
	p.EffectType = EffectType(p.EffectValue)

	// Read flags byte
	flags, err := r.ReadByte()
	if err != nil {
		return err
	}

	// Read TargetId if flag is set
	if (flags & EffectBitId) != 0 {
		p.TargetId, err = r.ReadCompressedInt()
		if err != nil {
			return err
		}
	}

	// Read PosA.X if flag is set
	if (flags & EffectBitPos1X) != 0 {
		x, err := r.ReadFloat32()
		if err != nil {
			return err
		}
		p.PosA.X = float64(x)
	}

	// Read PosA.Y if flag is set
	if (flags & EffectBitPos1Y) != 0 {
		y, err := r.ReadFloat32()
		if err != nil {
			return err
		}
		p.PosA.Y = float64(y)
	}

	// Read PosB.X if flag is set
	if (flags & EffectBitPos2X) != 0 {
		x, err := r.ReadFloat32()
		if err != nil {
			return err
		}
		p.PosB.X = float64(x)
	}

	// Read PosB.Y if flag is set
	if (flags & EffectBitPos2Y) != 0 {
		y, err := r.ReadFloat32()
		if err != nil {
			return err
		}
		p.PosB.Y = float64(y)
	}

	// Read Color if flag is set
	if (flags & EffectBitColor) != 0 {
		err = p.Color.Read(r)
		if err != nil {
			return err
		}
	}

	// Read Duration if flag is set
	if (flags & EffectBitDuration) != 0 {
		p.Duration, err = r.ReadFloat32()
		if err != nil {
			return err
		}
	}

	// Read UnknownValue if flag is set
	if (flags & UnknownBitId) != 0 {
		p.UnknownValue, err = r.ReadByte()
		if err != nil {
			return err
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *ShowEffect) Write(w interfaces.Writer) error {
	// Create a buffer to write the data
	var buffer bytes.Buffer
	bufferWriter := createBufferWriter(&buffer)

	// Calculate flags
	var flags byte = 0

	// Check if TargetId should be written
	if p.TargetId != 0 {
		flags |= EffectBitId
		if err := bufferWriter.WriteCompressedInt(p.TargetId); err != nil {
			return err
		}
	}

	// Check if PosA.X should be written
	if p.PosA.X != 0 {
		flags |= EffectBitPos1X
		if err := bufferWriter.WriteFloat32(float32(p.PosA.X)); err != nil {
			return err
		}
	}

	// Check if PosA.Y should be written
	if p.PosA.Y != 0 {
		flags |= EffectBitPos1Y
		if err := bufferWriter.WriteFloat32(float32(p.PosA.Y)); err != nil {
			return err
		}
	}

	// Check if PosB.X should be written
	if p.PosB.X != 0 {
		flags |= EffectBitPos2X
		if err := bufferWriter.WriteFloat32(float32(p.PosB.X)); err != nil {
			return err
		}
	}

	// Check if PosB.Y should be written
	if p.PosB.Y != 0 {
		flags |= EffectBitPos2Y
		if err := bufferWriter.WriteFloat32(float32(p.PosB.Y)); err != nil {
			return err
		}
	}

	// Check if Color should be written
	emptyColor := dataobjects.EmptyARGB()
	if !p.Color.Equals(emptyColor) {
		flags |= EffectBitColor
		if err := p.Color.Write(bufferWriter); err != nil {
			return err
		}
	}

	// Check if Duration should be written
	if p.Duration != 1.0 {
		flags |= EffectBitDuration
		if err := bufferWriter.WriteFloat32(p.Duration); err != nil {
			return err
		}
	}

	// Check if UnknownValue should be written
	if p.UnknownValue != 100 {
		flags |= UnknownBitId
		if err := bufferWriter.WriteByte(p.UnknownValue); err != nil {
			return err
		}
	}

	// Write EffectValue and flags to the main writer
	if err := w.WriteByte(p.EffectValue); err != nil {
		return err
	}
	if err := w.WriteByte(flags); err != nil {
		return err
	}

	// Write the buffer contents to the main writer
	if err := w.WriteBytes(buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

// createBufferWriter creates a writer for a buffer
func createBufferWriter(buffer *bytes.Buffer) interfaces.Writer {
	writer := packets.NewPacketWriter()
	return writer
}

// bufferWriter is a simple implementation of interfaces.Writer that writes to a buffer
type bufferWriter struct {
	buffer *bytes.Buffer
}

// Implement all the required methods of interfaces.Writer
func (bw *bufferWriter) WriteBytes(b []byte) error {
	_, err := bw.buffer.Write(b)
	return err
}

func (bw *bufferWriter) WriteByte(b byte) error {
	return bw.buffer.WriteByte(b)
}

func (bw *bufferWriter) WriteInt16(i int16) error {
	// Placeholder implementation
	return nil
}

func (bw *bufferWriter) WriteInt32(i int32) error {
	// Placeholder implementation
	return nil
}

func (bw *bufferWriter) WriteUInt16(i uint16) error {
	// Placeholder implementation
	return nil
}

func (bw *bufferWriter) WriteUInt32(i uint32) error {
	// Placeholder implementation
	return nil
}

func (bw *bufferWriter) WriteFloat32(f float32) error {
	// Placeholder implementation
	return nil
}

func (bw *bufferWriter) WriteString(s string) error {
	// Placeholder implementation
	return nil
}

func (bw *bufferWriter) WriteBool(b bool) error {
	// Placeholder implementation
	return nil
}

func (bw *bufferWriter) WriteCompressedInt(i int) error {
	// Placeholder implementation
	return nil
}

func (p *ShowEffect) ID() int32 {
	return int32(interfaces.ShowEffect)
}