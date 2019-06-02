package weekendraytracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

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

// NewCamera creates a camera with the given vertical field-of-view
// and horizontal-to-vertical aspect ratio.
func NewCamera(fov float64, aspect float64) Camera {
	theta := fov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := halfHeight * aspect
	return Camera{
		LowerLeftCorner: mgl64.Vec3{-halfWidth, -halfHeight, -1.0},
		Horizontal:      mgl64.Vec3{2 * halfWidth, 0.0, 0.0},
		Vertical:        mgl64.Vec3{0.0, 2 * halfHeight, 0.0},
		Origin:          mgl64.Vec3{0.0, 0.0, 0.0},
	}
}
