package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// SlotObject represents an item slot in the game
type SlotObject struct {
	ObjectID   int32
	SlotID     int32
	ObjectType int32
}

// NewSlotObject creates a new empty SlotObject
func NewSlotObject() *SlotObject {
	return &SlotObject{}
}

// NewSlotObjectWithData creates a new SlotObject with the given data
func NewSlotObjectWithData(objectID, slotID, objectType int32) *SlotObject {
	return &SlotObject{
		ObjectID:   objectID,
		SlotID:     slotID,
		ObjectType: objectType,
	}
}

// Read reads the slot object data from a Reader
func (s *SlotObject) Read(r interfaces.Reader) error {
	var err error
	s.ObjectID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	s.SlotID, err = r.ReadInt32()
	if err != nil {
		return err
	}
	s.ObjectType, err = r.ReadInt32()
	return err
}

// Write writes the slot object data to a PacketWriter
func (s *SlotObject) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(s.ObjectID); err != nil {
		return err
	}
	if err := w.WriteInt32(s.SlotID); err != nil {
		return err
	}
	return w.WriteInt32(s.ObjectType)
}

// Clone creates a copy of the SlotObject
func (s *SlotObject) Clone() *SlotObject {
	return &SlotObject{
		ObjectID:   s.ObjectID,
		ObjectType: s.ObjectType,
		SlotID:     s.SlotID,
	}
}

// String returns a string representation of the SlotObject
func (s *SlotObject) String() string {
	return fmt.Sprintf("{ ObjectId=%d, SlotId=%d, ObjectType=%d }", s.ObjectID, s.SlotID, s.ObjectType)
}
