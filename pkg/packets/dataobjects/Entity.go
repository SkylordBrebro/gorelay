package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// Entity represents a game entity with a type and status
type Entity struct {
	ObjectType uint16
	Status     *Status
}

// NewEntity creates a new Entity instance
func NewEntity() *Entity {
	return &Entity{
		Status: NewStatus(),
	}
}

// Read reads the entity data from a reader
func (e *Entity) Read(r interfaces.Reader) error {
	var err error
	e.ObjectType, err = r.ReadUInt16()
	if err != nil {
		return err
	}
	if e.Status == nil {
		e.Status = NewStatus()
	}
	return e.Status.Read(r)
}

// Write writes the entity data to a writer
func (e *Entity) Write(w interfaces.Writer) error {
	if err := w.WriteUInt16(e.ObjectType); err != nil {
		return err
	}
	if e.Status == nil {
		e.Status = NewStatus()
	}
	return e.Status.Write(w)
}

// Clone creates a copy of the Entity
func (e *Entity) Clone() DataObject {
	clone := NewEntity()
	clone.ObjectType = e.ObjectType
	if e.Status != nil {
		clone.Status = e.Status.Clone().(*Status)
	}
	return clone
}

// String returns a string representation of the Entity
func (e *Entity) String() string {
	status := "nil"
	if e.Status != nil {
		status = e.Status.String()
	}
	return fmt.Sprintf("{ ObjectType=%d, Status=%s }", e.ObjectType, status)
}
