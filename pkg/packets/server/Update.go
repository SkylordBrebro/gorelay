package server

import (
	"gorelay/pkg/packets/dataobjects"
	"gorelay/pkg/packets/interfaces"
)

// Update represents the server packet for updating game state
type Update struct {
	PlayerPosition *dataobjects.Location
	UnknownByte    byte
	Tiles          []*dataobjects.Tile
	NewObjs        []*dataobjects.Entity
	Drops          []int32
}

// Type returns the packet type for Update
func (p *Update) Type() interfaces.PacketType {
	return interfaces.Update
}

// Read reads the packet data from the provided reader
func (p *Update) Read(r interfaces.Reader) error {
	var err error

	// Read player position
	p.PlayerPosition = dataobjects.NewLocation()
	if err = p.PlayerPosition.Read(r); err != nil {
		// If we can't read position, create a default one
		p.PlayerPosition.X = 0
		p.PlayerPosition.Y = 0
	}

	// Read unknown byte
	p.UnknownByte, err = r.ReadByte()
	if err != nil {
		// If we can't read the byte, just set it to 0
		p.UnknownByte = 0
	}

	// Read tiles
	tileCount, err := r.ReadInt16()
	if err != nil {
		// If we can't read tile count, assume 0
		tileCount = 0
	}

	// Initialize tiles slice
	p.Tiles = make([]*dataobjects.Tile, 0)
	if tileCount > 0 && tileCount < 16384 { // Reasonable upper limit
		for i := int16(0); i < tileCount; i++ {
			tile := dataobjects.NewTile()
			if err := tile.Read(r); err != nil {
				// If we can't read a tile, stop reading tiles but continue with packet
				break
			}
			p.Tiles = append(p.Tiles, tile)
		}
	}

	// Read new objects
	newCount, err := r.ReadInt16()
	if err != nil {
		// If we can't read new object count, assume 0
		newCount = 0
	}

	// Initialize new objects slice
	p.NewObjs = make([]*dataobjects.Entity, 0)
	if newCount > 0 && newCount < 16384 { // Reasonable upper limit
		for i := int16(0); i < newCount; i++ {
			obj := dataobjects.NewEntity()
			if err := obj.Read(r); err != nil {
				// If we can't read an object, stop reading objects but continue with packet
				break
			}
			p.NewObjs = append(p.NewObjs, obj)
		}
	}

	// Read drops
	dropCount, err := r.ReadInt16()
	if err != nil {
		// If we can't read drop count, assume 0
		dropCount = 0
	}

	// Initialize drops slice
	p.Drops = make([]int32, 0)
	if dropCount > 0 && dropCount < 16384 { // Reasonable upper limit
		for i := int16(0); i < dropCount; i++ {
			drop, err := r.ReadInt32()
			if err != nil {
				// If we can't read a drop, stop reading drops
				break
			}
			p.Drops = append(p.Drops, drop)
		}
	}

	return nil
}

// Write writes the packet data to the provided writer
func (p *Update) Write(w interfaces.Writer) error {
	// Write player position
	if err := p.PlayerPosition.Write(w); err != nil {
		return err
	}

	// Write unknown byte
	if err := w.WriteByte(p.UnknownByte); err != nil {
		return err
	}

	// Write tiles
	if err := w.WriteInt16(int16(len(p.Tiles))); err != nil {
		return err
	}
	for _, tile := range p.Tiles {
		if err := tile.Write(w); err != nil {
			return err
		}
	}

	// Write new objects
	if err := w.WriteInt16(int16(len(p.NewObjs))); err != nil {
		return err
	}
	for _, obj := range p.NewObjs {
		if err := obj.Write(w); err != nil {
			return err
		}
	}

	// Write drops
	if err := w.WriteInt16(int16(len(p.Drops))); err != nil {
		return err
	}
	for _, drop := range p.Drops {
		if err := w.WriteInt32(drop); err != nil {
			return err
		}
	}

	return nil
}

func (p *Update) ID() int32 {
	return int32(interfaces.Update)
}
