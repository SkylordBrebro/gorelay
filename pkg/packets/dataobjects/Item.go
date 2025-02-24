package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// Item represents a game item
type Item struct {
	ItemItem int32
	SlotType int32
	Tradable bool
	Included bool
	ItemData string
}

// NewItem creates a new Item instance
func NewItem() *Item {
	return &Item{}
}

// Read reads the item data from a reader
func (i *Item) Read(r interfaces.Reader) error {
	var err error
	i.ItemItem, err = r.ReadInt32()
	if err != nil {
		return err
	}
	i.SlotType, err = r.ReadInt32()
	if err != nil {
		return err
	}
	i.Tradable, err = r.ReadBool()
	if err != nil {
		return err
	}
	i.Included, err = r.ReadBool()
	if err != nil {
		return err
	}
	i.ItemData, err = r.ReadString()
	return err
}

// Write writes the item data to a writer
func (i *Item) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(i.ItemItem); err != nil {
		return err
	}
	if err := w.WriteInt32(i.SlotType); err != nil {
		return err
	}
	if err := w.WriteBool(i.Tradable); err != nil {
		return err
	}
	if err := w.WriteBool(i.Included); err != nil {
		return err
	}
	return w.WriteString(i.ItemData)
}

// Clone creates a copy of the Item
func (i *Item) Clone() DataObject {
	return &Item{
		ItemItem: i.ItemItem,
		SlotType: i.SlotType,
		Tradable: i.Tradable,
		Included: i.Included,
		ItemData: i.ItemData,
	}
}

// String returns a string representation of the Item
func (i *Item) String() string {
	return fmt.Sprintf("{ ItemItem=%d, SlotType=%d, Tradable=%v, Included=%v, ItemData=%s }",
		i.ItemItem, i.SlotType, i.Tradable, i.Included, i.ItemData)
}
