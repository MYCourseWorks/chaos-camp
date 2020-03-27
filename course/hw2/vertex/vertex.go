package vertex

import "math"

// Vertex comment
type Vertex struct {
	X, Y float64
}

// Distance comment
func (v *Vertex) Distance(other *Vertex) float64 {
	return math.Hypot(other.X-v.X, other.Y-v.Y)
}

// Translate comment
func (v *Vertex) Translate(other *Vertex) {
	v.X += other.X
	v.Y += other.Y
}

// Rotate comment
func (v *Vertex) Rotate(rotationPoint *Vertex, angleInDegrees float64) {
	angleInRadians := angleInDegrees * (math.Pi / 180)
	cosTheta := math.Cos(angleInRadians)
	sinTheta := math.Sin(angleInRadians)

	v.X -= rotationPoint.X
	v.Y -= rotationPoint.Y

	newX := v.X*cosTheta - v.Y*sinTheta
	newY := v.X*sinTheta + v.Y*cosTheta

	v.X = newX + rotationPoint.X
	v.Y = newY + rotationPoint.Y
}

// IsEqualTo checks if the coordinates of two vertecies are the equivalent
func (v *Vertex) IsEqualTo(other *Vertex) bool {
	return math.Round(v.X) == math.Round(other.X) &&
		math.Round(v.Y) == math.Round(other.Y)
}
