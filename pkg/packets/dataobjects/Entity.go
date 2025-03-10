package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// Entity represents a game entity with its properties
type Entity struct {
	ObjectType int16
	Status     *Status
	Position   *Location
}

// NewEntity creates a new Entity instance
func NewEntity() *Entity {
	return &Entity{
		Status:   NewStatus(),
		Position: NewLocation(),
	}
}

// Read reads the entity data from the provided reader
func (e *Entity) Read(r interfaces.Reader) error {
	var err error

	// Read object type
	e.ObjectType, err = r.ReadInt16()
	if err != nil {
		return fmt.Errorf("failed to read object type: %v", err)
	}

	// Initialize position if needed
	if e.Position == nil {
		e.Position = NewLocation()
	}

	// Read position
	if err = e.Position.Read(r); err != nil {
		return fmt.Errorf("failed to read entity position: %v", err)
	}

	// Initialize status if needed
	if e.Status == nil {
		e.Status = NewStatus()
	}

	// Read status
	if err = e.Status.Read(r); err != nil {
		// If status read fails, just continue with empty status
		e.Status = NewStatus()
	}

	return nil
}

// Write writes the entity data to the provided writer
func (e *Entity) Write(w interfaces.Writer) error {
	if err := w.WriteInt16(e.ObjectType); err != nil {
		return err
	}

	if err := e.Position.Write(w); err != nil {
		return err
	}

	if err := e.Status.Write(w); err != nil {
		return err
	}

	return nil
}

// Clone creates a copy of the Entity
func (e *Entity) Clone() DataObject {
	clone := NewEntity()
	clone.ObjectType = e.ObjectType
	if e.Status != nil {
		clone.Status = e.Status.Clone().(*Status)
	}
	if e.Position != nil {
		clone.Position = e.Position.Clone().(*Location)
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
