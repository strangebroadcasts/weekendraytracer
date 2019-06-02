package weekendraytracer

import "github.com/go-gl/mathgl/mgl64"

// Camera is an axis-aligned camera.
type Camera struct {
	LowerLeftCorner mgl64.Vec3
	Horizontal      mgl64.Vec3
	Vertical        mgl64.Vec3
	Origin          mgl64.Vec3
}

// GetRay transforms screen-space coordinates of the camera
// to a world-space ray:
func (c Camera) GetRay(u float64, v float64) Ray {
	direction := c.LowerLeftCorner.Add(c.Horizontal.Mul(u).Add(c.Vertical.Mul(v)))
	return Ray{A: c.Origin, B: direction}
}
