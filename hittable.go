package weekendraytracer

import (
	"github.com/go-gl/mathgl/mgl64"
)

// HitRecord stores information about a ray intersection.
type HitRecord struct {
	// T is the parameter in the ray A + tB where the intersection is
	T float64
	// P is the point the intersection occurred at
	P mgl64.Vec3
	// Normal is the surface's normal at the intersection point
	Normal mgl64.Vec3
	// Mat is the material of the surface at the intersection point
	Mat Material
}

// Hittable objects can be tested for intersection against a ray.
type Hittable interface {
	// Hit lists the intersections of a ray against this primitive.
	Hit(r Ray, tMin float64, tMax float64) []HitRecord
	// Material returns the material this primitive uses.
	Material() Material
}

// HittableList is a list of primitives which can be tested for intersection.
type HittableList []Hittable

// Hit checks whether any element in the list intersects the ray,
// and returns the closest intersection (if any) in a single-element array.
func (l HittableList) Hit(r Ray, tMin float64, tMax float64) []HitRecord {
	recs := make([]HitRecord, 0, 1)
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
				bestIntersection.Mat = intersection.Mat
			}
		}
	}
	if hitAnything {
		recs = append(recs, bestIntersection)
	}
	return recs
}
