package weekendraytracer

import "github.com/go-gl/mathgl/mgl64"

// Ray is a parametric line.
type Ray struct {
	// A is the origin of the ray.
	A mgl64.Vec3
	// B is the direction vector of the ray.
	B mgl64.Vec3
}

// Origin returns the ray's origin
func (r Ray) Origin() mgl64.Vec3 {
	return r.A
}

// Direction returns the ray's direction
func (r Ray) Direction() mgl64.Vec3 {
	return r.B
}

// PointAtParameter gets a point from the ray at time t:
func (r Ray) PointAtParameter(t float64) mgl64.Vec3 {
	// p(t) = A + t*B
	return r.A.Add(r.B.Mul(t))
}
