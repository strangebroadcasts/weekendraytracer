package weekendraytracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

// HitRecord stores information about a ray intersection.
type HitRecord struct {
	T      float64
	P      mgl64.Vec3
	Normal mgl64.Vec3
}

// Hittable objects can be tested for intersection against a ray.
type Hittable interface {
	// Hit lists the intersections of a ray against this primitive.
	Hit(r Ray, tMin float64, tMax float64) []HitRecord
}

// HittableList is a list of primitives which can be tested for intersection.
type HittableList []Hittable

// Hit checks whether any element in the list intersects the ray,
// and returns the closest intersection (if any) in a single-element array.
func (l HittableList) Hit(r Ray, tMin float64, tMax float64) []HitRecord {
	recs := make([]HitRecord, 0)
	bestIntersection := HitRecord{T: tMax}
	hitAnything := false
	for _, primitive := range l {
		hits := primitive.Hit(r, tMin, tMax)
		for _, intersection := range hits {
			if intersection.T < bestIntersection.T {
				hitAnything = true
				bestIntersection.T = intersection.T
				bestIntersection.P = intersection.P
				bestIntersection.Normal = intersection.Normal
			}
		}
	}
	if hitAnything {
		recs = append(recs, bestIntersection)
	}
	return recs
}
