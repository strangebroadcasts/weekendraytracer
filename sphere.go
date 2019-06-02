package weekendraytracer

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// Sphere is a sphere defined by its center and radius:
type Sphere struct {
	Center mgl64.Vec3
	Radius float64
}

// Hit checks whether the ray r intersects this sphere:
func (s Sphere) Hit(r Ray, tMin float64, tMax float64) []HitRecord {
	recs := make([]HitRecord, 0, 2)
	oc := r.Origin().Sub(s.Center)
	a := r.Direction().Dot(r.Direction())
	b := oc.Dot(r.Direction())
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - a*c
	// If the discriminant is positive, we have one or
	// two intersections:
	if discriminant > 0 {
		temp := (-b - math.Sqrt(b*b-a*c)) / a
		if temp < tMax && temp > tMin {
			p := r.PointAtParameter(temp)
			normal := r.PointAtParameter(temp).Sub(s.Center).Mul(1.0 / s.Radius)
			hit := HitRecord{T: temp, P: p, Normal: normal}
			recs = append(recs, hit)
		}
		temp = (-b + math.Sqrt(b*b-a*c)) / a
		if temp < tMax && temp > tMin {
			p := r.PointAtParameter(temp)
			normal := r.PointAtParameter(temp).Sub(s.Center).Mul(1.0 / s.Radius)
			hit := HitRecord{T: temp, P: p, Normal: normal}
			recs = append(recs, hit)
		}
	}
	return recs
}
