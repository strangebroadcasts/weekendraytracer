package weekendraytracer

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

// Camera is an axis-aligned camera.
type Camera struct {
	LowerLeftCorner mgl64.Vec3
	Horizontal      mgl64.Vec3
	Vertical        mgl64.Vec3
	Origin          mgl64.Vec3

	u          mgl64.Vec3
	v          mgl64.Vec3
	w          mgl64.Vec3
	lensRadius float64
}

// randomInUnitDisk gets a random vector on the unit disk:
func randomInUnitDisk() mgl64.Vec3 {
	p := mgl64.Vec3{99.0, 99.0, 0.0}
	for p.Dot(p) >= 1.0 {
		x := (2.0 * rand.Float64()) - 1.0
		y := (2.0 * rand.Float64()) - 1.0
		p = mgl64.Vec3{x, y, 0.0}
	}
	return p
}

// GetRay transforms screen-space coordinates of the camera
// to a world-space ray:
func (c Camera) GetRay(u float64, v float64) Ray {
	direction := c.LowerLeftCorner.Add(c.Horizontal.Mul(u))
	direction = direction.Add(c.Vertical.Mul(v))
	direction = direction.Sub(c.Origin)
	rd := randomInUnitDisk().Mul(c.lensRadius)
	offset := c.u.Mul(rd.X()).Add(c.v.Mul(rd.Y()))
	return Ray{A: c.Origin.Add(offset), B: direction.Sub(offset)}
}

// NewCamera creates a camera with origin in (from),
// pointed at the point (at),
// using (up) as the world up vector,
// with the given vertical field-of-view
// and horizontal-to-vertical aspect ratio.
func NewCamera(from mgl64.Vec3, at mgl64.Vec3, up mgl64.Vec3, fov, aspect, aperture, focus float64) Camera {
	theta := fov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := halfHeight * aspect

	w := from.Sub(at).Normalize()
	u := up.Cross(w).Normalize()
	v := w.Cross(u)

	llc := from.Sub(u.Mul(halfWidth * focus))
	llc = llc.Sub(v.Mul(halfHeight * focus))
	llc = llc.Sub(w.Mul(focus))

	return Camera{
		LowerLeftCorner: llc,
		Horizontal:      u.Mul(2 * halfWidth * focus),
		Vertical:        v.Mul(2 * halfHeight * focus),
		Origin:          from,
		lensRadius:      aperture / 2.0,
		u:               u,
		v:               v,
		w:               w,
	}
}
