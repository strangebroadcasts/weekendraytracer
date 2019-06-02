package weekendraytracer

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
)

// Epsilon is the lower distance bound for reflections
// (used to combat shadow acne)
const Epsilon = 0.001

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

// Get the pixel color for this ray.
// (Called "color" in RTiaW, which conflicts with the color package)
func getColor(r Ray, world HittableList) mgl64.Vec3 {
	hits := world.Hit(r, Epsilon, math.MaxFloat64)
	if len(hits) > 0 {
		// HittableList makes sure the first (and only) intersection
		// is the closest one:
		hit := hits[0]
		// Find the ray reflecting off the surface:
		target := hit.P.Add(hit.Normal).Add(randomInsideUnitSphere())
		reflection := Ray{A: hit.P, B: target.Sub(hit.P)}
		// Our initial diffuse surface absorbs 50% of the light hitting it.
		return getColor(reflection, world).Mul(0.5)
	}
	// If we don't intersect with anything, plot a background instead:
	unitDirection := r.Direction().Normalize()
	t := 0.5*unitDirection.Y() + 1.0
	// Linearly interpolate between white and blue:
	A := mgl64.Vec3{1.0, 1.0, 1.0}.Mul(1.0 - t)
	B := mgl64.Vec3{0.5, 0.7, 1.0}.Mul(t)
	return A.Add(B)
}

// Render ray-traces the scene, outputting an image with the given dimensions.
func Render(width int, height int, samples int) image.Image {
	canvas := image.NewNRGBA(image.Rect(0, 0, width, height))

	cam := Camera{
		LowerLeftCorner: mgl64.Vec3{-2.0, -1.0, -1.0},
		Horizontal:      mgl64.Vec3{4.0, 0.0, 0.0},
		Vertical:        mgl64.Vec3{0.0, 2.0, 0.0},
		Origin:          mgl64.Vec3{0.0, 0.0, 0.0},
	}

	world := make(HittableList, 2)
	world[0] = Sphere{Center: mgl64.Vec3{0.0, 0.0, -1.0}, Radius: 0.5}
	world[1] = Sphere{Center: mgl64.Vec3{0.0, -100.5, -1.0}, Radius: 100}

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			col := mgl64.Vec3{}
			for s := 0; s < samples; s++ {
				u := (float64(i) + rand.Float64()) / float64(width)
				// Note that we flip the vertical axis here.
				v := (float64(height-j) + rand.Float64()) / float64(height)
				r := cam.GetRay(u, v)
				col = col.Add(getColor(r, world))
			}
			col = col.Mul(1.0 / float64(samples))
			// Gamma-correct the colors:
			col = mgl64.Vec3{math.Sqrt(col[0]), math.Sqrt(col[1]), math.Sqrt(col[2])}
			ir, ig, ib := uint8(255.99*col.X()), uint8(255.99*col.Y()), uint8(255.99*col.Z())
			canvas.Set(i, j, color.NRGBA{R: ir, G: ig, B: ib, A: 255})
		}
	}
	return canvas
}
