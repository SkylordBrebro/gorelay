package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// Status represents an entity's status including position and stats
type Status struct {
	ObjectID int
	Position *Location
	Data     []*StatData
}

// NewStatus creates a new Status instance
func NewStatus() *Status {
	return &Status{
		Position: NewLocation(),
		Data:     make([]*StatData, 0),
	}
}

// Read reads the status data from a Reader
func (s *Status) Read(r interfaces.Reader) error {
	var err error
	s.ObjectID, err = r.ReadCompressedInt()
	if err != nil {
		return err
	}

	if err = s.Position.Read(r); err != nil {
		return err
	}

	capacity, err := r.ReadCompressedInt()
	if err != nil {
		return err
	}

	s.Data = make([]*StatData, capacity)
	for i := 0; i < capacity; i++ {
		s.Data[i] = NewStatData()
		if err = s.Data[i].Read(r); err != nil {
			return err
		}
	}
	return nil
}

// Write writes the status data to a Writer
func (s *Status) Write(w interfaces.Writer) error {
	if err := w.WriteCompressedInt(s.ObjectID); err != nil {
		return err
	}

	if err := s.Position.Write(w); err != nil {
		return err
	}

	if err := w.WriteCompressedInt(len(s.Data)); err != nil {
		return err
	}

	for _, stat := range s.Data {
		if err := stat.Write(w); err != nil {
			return err
		}
	}
	return nil
}

// Clone creates a copy of the Status
func (s *Status) Clone() DataObject {
	newData := make([]*StatData, len(s.Data))
	for i, stat := range s.Data {
		if stat != nil {
			newData[i] = stat.Clone().(*StatData)
		}
	}

	clone := NewStatus()
	clone.ObjectID = s.ObjectID
	if s.Position != nil {
		clone.Position = s.Position.Clone().(*Location)
	}
	clone.Data = newData
	return clone
}

// String returns a string representation of the Status
func (s *Status) String() string {
	stats := make([]string, len(s.Data))
	for i, stat := range s.Data {
		if stat != nil {
			stats[i] = stat.String()
		}
	}
	return fmt.Sprintf("{ ObjectId=%d, Position=%v, Stats=[%s] }",
		s.ObjectID, s.Position, stats)
}
