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
		return err
	}

	// Read unknown byte
	p.UnknownByte, err = r.ReadByte()
	if err != nil {
		return err
	}

	// Read tiles
	tileCount, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.Tiles = make([]*dataobjects.Tile, tileCount)
	for i := 0; i < tileCount; i++ {
		p.Tiles[i] = dataobjects.NewTile()
		if err = p.Tiles[i].Read(r); err != nil {
			return err
		}
	}

	// Read new objects
	objCount, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.NewObjs = make([]*dataobjects.Entity, objCount)
	for i := 0; i < objCount; i++ {
		p.NewObjs[i] = dataobjects.NewEntity()
		if err = p.NewObjs[i].Read(r); err != nil {
			return err
		}
	}

	// Read drops
	dropCount, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}
	p.Drops = make([]int32, dropCount)
	for i := 0; i < dropCount; i++ {
		dropValue, err := r.ReadCompressedInt()
		if err != nil {
			return err
		}
		p.Drops[i] = int32(dropValue)
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
	if err := w.WriteCompressedInt(len(p.Tiles)); err != nil {
		return err
	}
	for _, tile := range p.Tiles {
		if err := tile.Write(w); err != nil {
			return err
		}
	}

	// Write new objects
	if err := w.WriteCompressedInt(len(p.NewObjs)); err != nil {
		return err
	}
	for _, obj := range p.NewObjs {
		if err := obj.Write(w); err != nil {
			return err
		}
	}

	// Write drops
	if err := w.WriteCompressedInt(len(p.Drops)); err != nil {
		return err
	}
	for _, drop := range p.Drops {
		if err := w.WriteCompressedInt(int(drop)); err != nil {
			return err
		}
	}

	return nil
}
