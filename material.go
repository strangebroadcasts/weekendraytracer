package weekendraytracer

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

// Material describes how light reflects off a surface
type Material interface {
	// Scatter describes the surface's response to the ray rIn,
	// at the intersection point given in hit.
	// Returns: whether the ray is scattered,
	// a vector describing the attenuation of the ray,
	// and the scattered ray (if any)
	Scatter(rIn Ray, hit HitRecord) (bool, mgl64.Vec3, Ray)
}

// Lambertian surfaces are diffuse, and scatter light in random directions.
type Lambertian struct {
	Albedo mgl64.Vec3
}

// Scatter simulates the reflection off a diffuse surface.
func (l Lambertian) Scatter(rIn Ray, hit HitRecord) (bool, mgl64.Vec3, Ray) {
	target := hit.P.Add(hit.Normal).Add(randomInsideUnitSphere())
	scattered := Ray{A: hit.P, B: target.Sub(hit.P)}
	attenuation := l.Albedo
	return true, attenuation, scattered
}

// randomInsideUnitSphere generates a random vector inside a unit sphere:
func randomInsideUnitSphere() mgl64.Vec3 {
	p := mgl64.Vec3{99.0, 99.0, 99.0}
	for p.LenSqr() >= 1.0 {
		x := (2.0 * rand.Float64()) - 1.0
		y := (2.0 * rand.Float64()) - 1.0
		z := (2.0 * rand.Float64()) - 1.0
		p = mgl64.Vec3{x, y, z}
	}
	return p
}

// Metallic surfaces reflect light.
type Metallic struct {
	Albedo    mgl64.Vec3
	Fuzziness float64
}

// Scatter simulates the reflection off a metallic surface.
func (m Metallic) Scatter(rIn Ray, hit HitRecord) (bool, mgl64.Vec3, Ray) {
	// To simulate uneven surfaces, we distort the reflection slightly:
	fuzz := randomInsideUnitSphere().Mul(m.Fuzziness)
	reflected := reflect(rIn.Direction().Normalize(), hit.Normal)
	scattered := Ray{A: hit.P, B: reflected.Add(fuzz)}
	attenuation := m.Albedo
	isScattered := scattered.Direction().Dot(hit.Normal) > 0
	return isScattered, attenuation, scattered
}

// reflect calculates the reflection of vector v against the surface normal n.
func reflect(v, n mgl64.Vec3) mgl64.Vec3 {
	t := n.Mul(2 * v.Dot(n))
	return v.Sub(t)
}
