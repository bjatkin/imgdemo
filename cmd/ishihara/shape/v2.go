package shape

import (
	"math"
	"math/rand"
)

// V2 is a vector2 consisting of two float64s
type V2 struct {
	X float64
	Y float64
}

// NewV2 creates a new V2
func NewV2(x, y float64) V2 {
	return V2{X: x, Y: y}
}

// V2Distance calculates the distance between V2 points in space
func V2Distance(v1, v2 V2) float64 {
	deltaX := (v2.X - v1.X)
	deltaY := (v2.Y - v1.Y)
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

// randomPolar creates a random polar coordinate with the given maximum radius
func randomPolar(maxRadius float64) V2 {
	angle := rand.Float64() * 2.0 * math.Pi
	radius := rand.Float64() * maxRadius
	return NewV2(radius, angle)
}

// polarToCartesian converts a polar coordinate to a cartesian one
func polarToCartesian(bounds *Circle, polar V2) V2 {
	x := polar.X * math.Cos(polar.Y)
	y := polar.X * math.Sin(polar.Y)
	return NewV2(x+bounds.Center.X, y+bounds.Center.Y)
}
