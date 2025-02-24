package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// StatData represents a single stat value
type StatData struct {
	ID          StatsType
	IntValue    int
	StringValue string
	ExtraValue  byte
}

// NewStatData creates a new StatData instance
func NewStatData() *StatData {
	return &StatData{}
}

// IsStringData returns true if the stat contains string data
func (s *StatData) IsStringData() bool {
	return s.ID.IsUTF()
}

// Read reads the stat data from a Reader
func (s *StatData) Read(r interfaces.Reader) error {
	statType, err := r.ReadByte()
	if err != nil {
		return err
	}
	s.ID = StatsType(statType)

	if s.IsStringData() {
		s.StringValue, err = r.ReadString()
		if err != nil {
			return err
		}
	} else {
		s.IntValue, err = r.ReadCompressedInt()
		if err != nil {
			return err
		}
	}

	s.ExtraValue, err = r.ReadByte()
	return err
}

// Write writes the stat data to a Writer
func (s *StatData) Write(w interfaces.Writer) error {
	if err := w.WriteByte(byte(s.ID)); err != nil {
		return err
	}

	if s.IsStringData() {
		if err := w.WriteString(s.StringValue); err != nil {
			return err
		}
	} else {
		if err := w.WriteCompressedInt(s.IntValue); err != nil {
			return err
		}
	}

	return w.WriteByte(s.ExtraValue)
}

// Clone creates a copy of the StatData
func (s *StatData) Clone() DataObject {
	return &StatData{
		ID:          s.ID,
		IntValue:    s.IntValue,
		StringValue: s.StringValue,
		ExtraValue:  s.ExtraValue,
	}
}

// String returns a string representation of the StatData
func (s *StatData) String() string {
	value := s.IntValue
	if s.IsStringData() {
		value = -1 // Indicate string value
	}
	return fmt.Sprintf("(%s = %v (Extra: %d))", s.ID, value, s.ExtraValue)
}
