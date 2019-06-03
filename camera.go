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
	direction := c.LowerLeftCorner.Add(c.Horizontal.Mul(u).Add(c.Vertical.Mul(v))).Sub(c.Origin)
	return Ray{A: c.Origin, B: direction}
}

// NewCamera creates a camera with origin in (from),
// pointed at the point (at),
// using (up) as the world up vecotr
// with the given vertical field-of-view
// and horizontal-to-vertical aspect ratio.
func NewCamera(from mgl64.Vec3, at mgl64.Vec3, up mgl64.Vec3, fov float64, aspect float64) Camera {
	theta := fov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := halfHeight * aspect

	w := from.Sub(at).Normalize()
	u := up.Cross(w).Normalize()
	v := w.Cross(u)

	llc := from.Sub(u.Mul(halfWidth)).Sub(v.Mul(halfHeight)).Sub(w)

	return Camera{
		LowerLeftCorner: llc,
		Horizontal:      u.Mul(2 * halfWidth),
		Vertical:        v.Mul(2 * halfHeight),
		Origin:          from,
	}
}
