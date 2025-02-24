package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// Tile represents a game tile with position and type
type Tile struct {
	X    int16
	Y    int16
	Type uint16
}

// NewTile creates a new empty Tile instance
func NewTile() *Tile {
	return &Tile{}
}

// NewTileWithData creates a new Tile with the given data
func NewTileWithData(x, y int16, tileType uint16) *Tile {
	return &Tile{
		X:    x,
		Y:    y,
		Type: tileType,
	}
}

// Read reads the tile data from a Reader
func (t *Tile) Read(r interfaces.Reader) error {
	var err error
	t.X, err = r.ReadInt16()
	if err != nil {
		return err
	}
	t.Y, err = r.ReadInt16()
	if err != nil {
		return err
	}
	t.Type, err = r.ReadUInt16()
	return err
}

// Write writes the tile data to a Writer
func (t *Tile) Write(w interfaces.Writer) error {
	if err := w.WriteInt16(t.X); err != nil {
		return err
	}
	if err := w.WriteInt16(t.Y); err != nil {
		return err
	}
	return w.WriteUInt16(t.Type)
}

// Clone creates a copy of the Tile
func (t *Tile) Clone() DataObject {
	return NewTileWithData(t.X, t.Y, t.Type)
}

// String returns a string representation of the Tile
func (t *Tile) String() string {
	return fmt.Sprintf("{ X=%d, Y=%d, Type=%d }", t.X, t.Y, t.Type)
}
