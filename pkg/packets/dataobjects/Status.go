package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// Status represents an entity's status data
type Status struct {
	ObjectID int32
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

// Read reads the status data from the provided reader
func (s *Status) Read(r interfaces.Reader) error {
	var err error

	// Read object ID
	s.ObjectID, err = r.ReadInt32()
	if err != nil {
		// If we can't read object ID, use 0
		s.ObjectID = 0
	}

	// Read position
	if s.Position == nil {
		s.Position = NewLocation()
	}
	if err = s.Position.Read(r); err != nil {
		// If we can't read position, use default values
		s.Position.X = 0
		s.Position.Y = 0
	}

	// Read stat data count
	statCount, err := r.ReadInt16()
	if err != nil {
		// If we can't read stat count, assume 0
		statCount = 0
	}

	// Initialize stat data slice
	s.Data = make([]*StatData, 0)
	if statCount > 0 && statCount < 128 { // Reasonable limit for stats
		for i := int16(0); i < statCount; i++ {
			stat := NewStatData()
			if err := stat.Read(r); err != nil {
				// If we can't read a stat, stop reading stats but continue with packet
				break
			}
			s.Data = append(s.Data, stat)
		}
	}

	return nil
}

// Write writes the status data to the provided writer
func (s *Status) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(s.ObjectID); err != nil {
		return err
	}

	if err := s.Position.Write(w); err != nil {
		return err
	}

	if err := w.WriteInt16(int16(len(s.Data))); err != nil {
		return err
	}

	for _, stat := range s.Data {
		if err := stat.Write(w); err != nil {
			return err
		}
	}

	return nil
}

// Clone creates a deep copy of the Status
func (s *Status) Clone() interface{} {
	clone := NewStatus()
	clone.ObjectID = s.ObjectID
	if s.Position != nil {
		clone.Position = s.Position.Clone().(*Location)
	}
	if s.Data != nil {
		clone.Data = make([]*StatData, len(s.Data))
		for i, stat := range s.Data {
			if stat != nil {
				clone.Data[i] = stat.Clone().(*StatData)
			}
		}
	}
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
