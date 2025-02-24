package services

import (
	"math/rand"
	"time"
)

// Random provides utility functions for random number generation
type Random struct {
	rng *rand.Rand
}

// NewRandom creates a new random number generator
func NewRandom() *Random {
	return &Random{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NextInt returns a random integer between min and max (inclusive)
func (r *Random) NextInt(min, max int) int {
	return min + r.rng.Intn(max-min+1)
}

// NextFloat returns a random float between min and max
func (r *Random) NextFloat(min, max float64) float64 {
	return min + r.rng.Float64()*(max-min)
}

// NextBool returns a random boolean value
func (r *Random) NextBool() bool {
	return r.rng.Float64() < 0.5
}

// NextIntInRange returns a random integer in the given range
func (r *Random) NextIntInRange(min, max int32) int32 {
	return min + int32(r.rng.Int31n(max-min+1))
}

// NextFloatInRange returns a random float in the given range
func (r *Random) NextFloatInRange(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
