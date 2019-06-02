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

// MaxBounces is how many light bounces we allow at most.
const MaxBounces = 10

// Get the pixel color for this ray.
// (Called "color" in RTiaW, which conflicts with the color package)
func getColor(r Ray, world HittableList, depth int) mgl64.Vec3 {
	hits := world.Hit(r, Epsilon, math.MaxFloat64)
	if len(hits) > 0 {
		// HittableList makes sure the first (and only) intersection
		// is the closest one:
		hit := hits[0]
		// Simulate the light response of this material.
		isScattered, attenuation, scattered := hit.Mat.Scatter(r, hit)
		// If this ray is reflected off the surface, determine the response
		// of the scattered ray:
		if isScattered && depth < MaxBounces {
			response := getColor(scattered, world, depth+1)

			return mgl64.Vec3{response.X() * attenuation.X(),
				response.Y() * attenuation.Y(),
				response.Z() * attenuation.Z()}
		}
		// Otherwise, the ray was absorbed, or we have exceeded the
		// maximum number of light bounces:
		return mgl64.Vec3{0.0, 0.0, 0.0}
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

	world := make(HittableList, 4)
	world[0] = Sphere{
		Center: mgl64.Vec3{0.0, 0.0, -1.0},
		Radius: 0.5,
		Mat:    Lambertian{Albedo: mgl64.Vec3{0.8, 0.3, 0.3}}}
	world[1] = Sphere{
		Center: mgl64.Vec3{0.0, -100.5, -1.0},
		Radius: 100,
		Mat:    Lambertian{Albedo: mgl64.Vec3{0.8, 0.8, 0.0}}}
	world[2] = Sphere{
		Center: mgl64.Vec3{1.0, 0.0, -1.0},
		Radius: 0.5,
		Mat:    Metallic{Albedo: mgl64.Vec3{0.8, 0.6, 0.2}, Fuzziness: 1.0}}
	world[3] = Sphere{
		Center: mgl64.Vec3{-1.0, 0.0, -1.0},
		Radius: 0.5,
		Mat:    Dielectric{Refractance: 1.5}}

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			col := mgl64.Vec3{}
			for s := 0; s < samples; s++ {
				u := (float64(i) + rand.Float64()) / float64(width)
				// Note that we flip the vertical axis here.
				v := (float64(height-j) + rand.Float64()) / float64(height)
				r := cam.GetRay(u, v)
				col = col.Add(getColor(r, world, 0))
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
