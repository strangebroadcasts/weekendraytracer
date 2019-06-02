package weekendraytracer

import (
	"math"
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

// Dielectric surfaces both reflect and refract light:
type Dielectric struct {
	Refractance float64
}

// refract simulates the refraction of vector v against the surface normal n,
// returning whether the ray is refracted (and the refracted ray if it is)
func refract(v mgl64.Vec3, n mgl64.Vec3, refractRatio float64) (bool, mgl64.Vec3) {
	uv := v.Normalize()
	dt := uv.Dot(n)
	discriminant := 1.0 - refractRatio*refractRatio*(1.0-dt*dt)
	if discriminant > 0.0 {
		a := uv.Sub(n.Mul(dt)).Mul(refractRatio)
		b := n.Mul(math.Sqrt(discriminant))
		refracted := a.Sub(b)
		return true, refracted
	}
	return false, mgl64.Vec3{}
}

// schlick is Christophe Schlick's polynomial approximation of
// the specular reflection coefficient:
func schlick(cosine float64, refract float64) float64 {
	r0 := (1.0 - refract) / (1.0 + refract)
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow((1-cosine), 5.0)
}

// Scatter simulates the reflection off a dielectric surface.
func (d Dielectric) Scatter(rIn Ray, hit HitRecord) (bool, mgl64.Vec3, Ray) {
	var outwardNormal mgl64.Vec3
	var refractRatio float64
	var cosine float64
	var reflectProb float64

	reflected := reflect(rIn.Direction(), hit.Normal)
	attenuation := mgl64.Vec3{1.0, 1.0, 1.0}
	if rIn.Direction().Dot(hit.Normal) > 0 {
		outwardNormal = hit.Normal.Mul(-1.0)
		refractRatio = d.Refractance
		cosine = d.Refractance * rIn.Direction().Dot(hit.Normal) / rIn.Direction().Len()
	} else {
		outwardNormal = hit.Normal
		refractRatio = 1.0 / d.Refractance
		cosine = -1.0 * rIn.Direction().Dot(hit.Normal) / rIn.Direction().Len()
	}

	isRefracted, refracted := refract(rIn.Direction(), outwardNormal, refractRatio)
	if isRefracted {
		reflectProb = schlick(cosine, d.Refractance)
	} else {
		reflectProb = 1.0
	}

	var scattered mgl64.Vec3
	if rand.Float64() < reflectProb {
		scattered = reflected
	} else {
		scattered = refracted
	}
	return true, attenuation, Ray{A: hit.P, B: scattered}
}
