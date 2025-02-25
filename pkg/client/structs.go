package client

import (
	"fmt"
	"math"
)

// Common direction constants
var (
	Up    = &WorldPosData{X: 0, Y: -1}
	Down  = &WorldPosData{X: 0, Y: 1}
	Left  = &WorldPosData{X: -1, Y: 0}
	Right = &WorldPosData{X: 1, Y: 0}
	Zero  = &WorldPosData{X: 0, Y: 0}
)

// WorldPosData represents a position in the game world
type WorldPosData struct {
	X float32
	Y float32
}

// SquareDistanceTo calculates the squared distance to another position
func (w *WorldPosData) SquareDistanceTo(other *WorldPosData) float32 {
	dx := w.X - other.X
	dy := w.Y - other.Y
	return dx*dx + dy*dy
}

// DistanceTo calculates the Euclidean distance to another position
func (w *WorldPosData) DistanceTo(other *WorldPosData) float32 {
	return float32(math.Sqrt(float64(w.SquareDistanceTo(other))))
}

// Add adds another position vector to this one
func (w *WorldPosData) Add(other *WorldPosData) *WorldPosData {
	return &WorldPosData{
		X: w.X + other.X,
		Y: w.Y + other.Y,
	}
}

// Subtract subtracts another position vector from this one
func (w *WorldPosData) Subtract(other *WorldPosData) *WorldPosData {
	return &WorldPosData{
		X: w.X - other.X,
		Y: w.Y - other.Y,
	}
}

// Scale multiplies the position by a scalar value
func (w *WorldPosData) Scale(factor float32) *WorldPosData {
	return &WorldPosData{
		X: w.X * factor,
		Y: w.Y * factor,
	}
}

// Normalize returns a unit vector in the same direction
func (w *WorldPosData) Normalize() *WorldPosData {
	length := w.DistanceTo(&WorldPosData{X: 0, Y: 0})
	if length == 0 {
		return &WorldPosData{X: 0, Y: 0}
	}
	return &WorldPosData{
		X: w.X / length,
		Y: w.Y / length,
	}
}

// AngleTo calculates the angle in radians to another position
func (w *WorldPosData) AngleTo(other *WorldPosData) float32 {
	return float32(math.Atan2(float64(other.Y-w.Y), float64(other.X-w.X)))
}

// MoveTowards moves this position towards a target by a specified distance
func (w *WorldPosData) MoveTowards(target *WorldPosData, distance float32) *WorldPosData {
	if w.SquareDistanceTo(target) <= distance*distance {
		return &WorldPosData{X: target.X, Y: target.Y}
	}

	angle := w.AngleTo(target)
	return &WorldPosData{
		X: w.X + float32(math.Cos(float64(angle)))*distance,
		Y: w.Y + float32(math.Sin(float64(angle)))*distance,
	}
}

// Clone creates a copy of this position
func (w *WorldPosData) Clone() *WorldPosData {
	return &WorldPosData{X: w.X, Y: w.Y}
}

// IsZero checks if this is a zero position (0,0)
func (w *WorldPosData) IsZero() bool {
	return w.X == 0 && w.Y == 0
}

// Lerp linearly interpolates between this position and another by t (0.0 to 1.0)
func (w *WorldPosData) Lerp(other *WorldPosData, t float32) *WorldPosData {
	if t <= 0 {
		return w.Clone()
	}
	if t >= 1 {
		return other.Clone()
	}
	return &WorldPosData{
		X: w.X + (other.X-w.X)*t,
		Y: w.Y + (other.Y-w.Y)*t,
	}
}

// ToGridPosition converts world position to grid coordinates (tile position)
func (w *WorldPosData) ToGridPosition() (int, int) {
	return int(math.Floor(float64(w.X))), int(math.Floor(float64(w.Y)))
}

// FromGridPosition creates a world position from grid coordinates (centered in tile)
func FromGridPosition(gridX, gridY int) *WorldPosData {
	return &WorldPosData{
		X: float32(gridX) + 0.5,
		Y: float32(gridY) + 0.5,
	}
}

// Rotate rotates the position around the origin by the given angle in radians
func (w *WorldPosData) Rotate(angleRadians float32) *WorldPosData {
	sin := float32(math.Sin(float64(angleRadians)))
	cos := float32(math.Cos(float64(angleRadians)))

	return &WorldPosData{
		X: w.X*cos - w.Y*sin,
		Y: w.X*sin + w.Y*cos,
	}
}

// RotateAround rotates the position around a given center point by the given angle in radians
func (w *WorldPosData) RotateAround(center *WorldPosData, angleRadians float32) *WorldPosData {
	// Translate to origin
	translated := w.Subtract(center)

	// Rotate
	rotated := translated.Rotate(angleRadians)

	// Translate back
	return rotated.Add(center)
}

// Magnitude returns the length of this position vector from origin
func (w *WorldPosData) Magnitude() float32 {
	return float32(math.Sqrt(float64(w.X*w.X + w.Y*w.Y)))
}

// SqrMagnitude returns the squared length of this position vector from origin
func (w *WorldPosData) SqrMagnitude() float32 {
	return w.X*w.X + w.Y*w.Y
}

// Dot returns the dot product of this position and another
func (w *WorldPosData) Dot(other *WorldPosData) float32 {
	return w.X*other.X + w.Y*other.Y
}

// String returns a string representation of the position
func (w *WorldPosData) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", w.X, w.Y)
}

// Equals checks if this position is equal to another within a small epsilon
func (w *WorldPosData) Equals(other *WorldPosData) bool {
	const epsilon = 0.001
	return math.Abs(float64(w.X-other.X)) < epsilon &&
		math.Abs(float64(w.Y-other.Y)) < epsilon
}

// ClampMagnitude limits the magnitude of this vector to the specified maximum
func (w *WorldPosData) ClampMagnitude(maxLength float32) *WorldPosData {
	sqrMag := w.SqrMagnitude()
	if sqrMag > maxLength*maxLength {
		return w.Normalize().Scale(maxLength)
	}
	return w.Clone()
}

// Distance calculates the distance between two positions (static utility function)
func Distance(a, b *WorldPosData) float32 {
	return a.DistanceTo(b)
}

// Lerp linearly interpolates between two positions (static utility function)
func Lerp(a, b *WorldPosData, t float32) *WorldPosData {
	return a.Lerp(b, t)
}

// DirectionFromAngle creates a normalized direction vector from an angle in radians
func DirectionFromAngle(angleRadians float32) *WorldPosData {
	return &WorldPosData{
		X: float32(math.Cos(float64(angleRadians))),
		Y: float32(math.Sin(float64(angleRadians))),
	}
}

// AngleBetween calculates the smallest angle between two direction vectors
func AngleBetween(from, to *WorldPosData) float32 {
	// Normalize the vectors
	fromNorm := from.Normalize()
	toNorm := to.Normalize()

	// Calculate the dot product
	dot := fromNorm.Dot(toNorm)

	// Clamp to avoid floating point errors
	if dot > 1.0 {
		dot = 1.0
	} else if dot < -1.0 {
		dot = -1.0
	}

	return float32(math.Acos(float64(dot)))
}

// PlayerData represents the player's character data
type PlayerData struct {
	// Basic stats
	Name         string
	Level        int32
	Exp          int32
	NextLevelExp int32

	// Health and mana
	HP    int32
	MaxHP int32
	MP    int32
	MaxMP int32

	// Fame stats
	Fame        int32
	CurrentFame int32
	Stars       int32

	// Account info
	AccountID string

	// Guild info
	GuildName string
	GuildRank int32

	// Inventory and consumables
	BackpackSlots int32
	Inventory     []int32
	Potions       []PotionData

	// Character stats
	Stats map[string]int32
}

type PotionData struct {
	id       int32
	quantity int8
}

// IsDead returns whether the enemy is dead
func (e *Enemy) IsDead() bool {
	return e.HP <= 0
}
