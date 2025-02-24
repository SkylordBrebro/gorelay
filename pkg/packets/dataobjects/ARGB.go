package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// ARGB represents a color with alpha, red, green, and blue components
type ARGB struct {
	A uint8
	R uint8
	G uint8
	B uint8
}

// NewARGB creates a new ARGB instance
func NewARGB() *ARGB {
	return &ARGB{}
}

// NewARGBWithValues creates a new ARGB with the given values
func NewARGBWithValues(a, r, g, b uint8) *ARGB {
	return &ARGB{
		A: a,
		R: r,
		G: g,
		B: b,
	}
}

// NewARGBFromUint creates a new ARGB from a uint32 value
func NewARGBFromUint(value uint32) *ARGB {
	return &ARGB{
		A: uint8((value >> 24) & 0xFF),
		R: uint8((value >> 16) & 0xFF),
		G: uint8((value >> 8) & 0xFF),
		B: uint8(value & 0xFF),
	}
}

// EmptyARGB returns an ARGB with maximum alpha (fully opaque)
func EmptyARGB() *ARGB {
	return &ARGB{A: 255, R: 0, G: 0, B: 0}
}

// Read reads the ARGB data from a Reader
func (a *ARGB) Read(r interfaces.Reader) error {
	var err error

	// Read the 4 bytes as a uint32
	value, err := r.ReadUInt32()
	if err != nil {
		return err
	}

	// Extract components
	a.A = uint8((value >> 24) & 0xFF)
	a.R = uint8((value >> 16) & 0xFF)
	a.G = uint8((value >> 8) & 0xFF)
	a.B = uint8(value & 0xFF)

	return nil
}

// Write writes the ARGB data to a Writer
func (a *ARGB) Write(w interfaces.Writer) error {
	// Combine components into a uint32
	value := uint32(a.A)<<24 | uint32(a.R)<<16 | uint32(a.G)<<8 | uint32(a.B)
	return w.WriteUInt32(value)
}

// ToUint32 converts the ARGB to a uint32 value
func (a *ARGB) ToUint32() uint32 {
	return uint32(a.A)<<24 | uint32(a.R)<<16 | uint32(a.G)<<8 | uint32(a.B)
}

// Clone creates a copy of the ARGB
func (a *ARGB) Clone() DataObject {
	return &ARGB{
		A: a.A,
		R: a.R,
		G: a.G,
		B: a.B,
	}
}

// String returns a string representation of the ARGB
func (a *ARGB) String() string {
	return fmt.Sprintf("{ A=%d, R=%d, G=%d, B=%d }", a.A, a.R, a.G, a.B)
}

// Equals checks if this ARGB equals another
func (a *ARGB) Equals(other *ARGB) bool {
	return a.A == other.A && a.R == other.R && a.G == other.G && a.B == other.B
}
