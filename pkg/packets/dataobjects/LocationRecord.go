package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
)

// LocationRecord represents a recorded location with a timestamp
type LocationRecord struct {
	Time     int32
	Position *Location
}

// NewLocationRecord creates a new LocationRecord instance
func NewLocationRecord() *LocationRecord {
	return &LocationRecord{
		Position: NewLocation(),
	}
}

// Read reads the location record from a PacketReader
func (l *LocationRecord) Read(r interfaces.Reader) error {
	var err error
	l.Time, err = r.ReadInt32()
	if err != nil {
		return err
	}

	if l.Position == nil {
		l.Position = NewLocation()
	}
	return l.Position.Read(r)
}

// Write writes the location record to a PacketWriter
func (l *LocationRecord) Write(w interfaces.Writer) error {
	if err := w.WriteInt32(l.Time); err != nil {
		return err
	}
	return l.Position.Write(w)
}

// Clone creates a copy of the LocationRecord
func (l *LocationRecord) Clone() DataObject {
	return &LocationRecord{
		Time:     l.Time,
		Position: l.Position.Clone().(*Location),
	}
}

// String returns a string representation of the LocationRecord
func (l *LocationRecord) String() string {
	return fmt.Sprintf("{ Time=%d, Position=%v }", l.Time, l.Position)
}
