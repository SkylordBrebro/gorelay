package util

import "math"

// SquareDistance calculates the square of the distance between two points
func SquareDistance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}

// Distance calculates the distance between two points
func Distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(SquareDistance(x1, y1, x2, y2))
}

// AngleTo calculates the angle between two points
func AngleTo(x1, y1, x2, y2 float64) float64 {
	return math.Atan2(y2-y1, x2-x1)
}

// Lerp performs linear interpolation between two values
func Lerp(start, end, t float64) float64 {
	return start + t*(end-start)
}

// Clamp constrains a value between min and max
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
