package dataobjects

import (
	"fmt"
	"gorelay/pkg/packets/interfaces"
	"math"
)

// Location represents a 2D position in the game
type Location struct {
	X float64
	Y float64
}

// NewLocation creates a new Location instance
func NewLocation() *Location {
	return &Location{}
}

// NewLocationWithCoords creates a new Location with the given coordinates
func NewLocationWithCoords(x, y float64) *Location {
	return &Location{X: x, Y: y}
}

// Empty returns a Location at the origin (0,0)
func Empty() *Location {
	return &Location{}
}

// Read reads the location data from a Reader
func (l *Location) Read(r interfaces.Reader) error {
	var err error
	x, err := r.ReadFloat32()
	if err != nil {
		return err
	}
	l.X = float64(x)

	y, err := r.ReadFloat32()
	if err != nil {
		return err
	}
	l.Y = float64(y)
	return nil
}

// Write writes the location data to a Writer
func (l *Location) Write(w interfaces.Writer) error {
	if err := w.WriteFloat32(float32(l.X)); err != nil {
		return err
	}
	return w.WriteFloat32(float32(l.Y))
}

// DistanceSquaredTo returns the squared distance to another location
func (l *Location) DistanceSquaredTo(other *Location) float64 {
	dx := other.X - l.X
	dy := other.Y - l.Y
	return dx*dx + dy*dy
}

// DistanceTo returns the distance to another location
func (l *Location) DistanceTo(other *Location) float64 {
	return math.Sqrt(l.DistanceSquaredTo(other))
}

// GetAngle returns the angle between two locations
func (l *Location) GetAngle(l1, l2 *Location) float64 {
	return math.Atan2(l2.Y-l1.Y, l2.X-l1.X)
}

// GetAngleFromCoords returns the angle between two coordinate pairs
func (l *Location) GetAngleFromCoords(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y2-y1, x2-x1)
}

// PositionInDirection returns a new location at a distance in a direction
func (l *Location) PositionInDirection(angle, distance float64) *Location {
	return &Location{
		X: l.X + math.Cos(angle)*distance,
		Y: l.Y + math.Sin(angle)*distance,
	}
}

// PositionInDirectionFromCoords returns a new location at a distance in a direction from coordinates
func PositionInDirectionFromCoords(sourceX, sourceY, angle, distance float64) *Location {
	return &Location{
		X: sourceX + math.Cos(angle)*distance,
		Y: sourceY + math.Sin(angle)*distance,
	}
}

// Add returns a new location that is the sum of this location and another
func (l *Location) Add(other *Location) *Location {
	return &Location{X: l.X + other.X, Y: l.Y + other.Y}
}

// Subtract returns a new location that is this location minus another
func (l *Location) Subtract(other *Location) *Location {
	return &Location{X: l.X - other.X, Y: l.Y - other.Y}
}

// ScaleBy multiplies the coordinates by a value
func (l *Location) ScaleBy(value float64) {
	l.X *= value
	l.Y *= value
}

// DotProduct returns the dot product with another location
func (l *Location) DotProduct(other *Location) float64 {
	return l.X*other.X + l.Y*other.Y
}

// Clone creates a copy of the Location
func (l *Location) Clone() DataObject {
	return &Location{
		X: l.X,
		Y: l.Y,
	}
}

// String returns a string representation of the Location
func (l *Location) String() string {
	return fmt.Sprintf("{ X=%f, Y=%f }", l.X, l.Y)
}

// StringShort returns a shorter string representation of the Location
func (l *Location) StringShort() string {
	return fmt.Sprintf("{ X=%.2f, Y=%.2f }", l.X, l.Y)
}
