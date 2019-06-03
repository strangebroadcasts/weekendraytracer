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

// CreateRandomScene creates the "floor with random spheres"
// scene on pg. 40 of RTiaW.
func CreateRandomScene() HittableList {
	n := 500
	scene := make(HittableList, 0, n)

	floor := Sphere{
		Center: mgl64.Vec3{0.0, -1000.0, 0},
		Radius: 1000,
		Mat:    Lambertian{Albedo: mgl64.Vec3{0.5, 0.5, 0.5}}}
	scene = append(scene, floor)

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			center := mgl64.Vec3{float64(a) + 0.9*rand.Float64(), 0.2, float64(b) + rand.Float64()}
			if center.Sub(mgl64.Vec3{4, 0.2, 0}).Len() <= 0.9 {
				continue
			}

			sphere := Sphere{Center: center, Radius: 0.2}
			materialChoice := rand.Float64()
			if materialChoice < 0.8 { // Choose diffuse material for 80% of spheres.
				r := rand.Float64() * rand.Float64()
				g := rand.Float64() * rand.Float64()
				b := rand.Float64() * rand.Float64()
				sphere.Mat = Lambertian{Albedo: mgl64.Vec3{r, g, b}}
			} else if materialChoice < 0.95 { // Let 15% of spheres be metallic:
				r := 0.5 * (1.0 + rand.Float64())
				g := 0.5 * (1.0 + rand.Float64())
				b := 0.5 * (1.0 + rand.Float64())
				fuzz := rand.Float64()
				sphere.Mat = Metallic{Albedo: mgl64.Vec3{r, g, b}, Fuzziness: fuzz}
			} else { // Let the rest be dielectric:
				sphere.Mat = Dielectric{Refractance: 1.5}
			}
			scene = append(scene, sphere)
		}
	}

	dielectricSphere := Sphere{
		Center: mgl64.Vec3{0, 1, 0},
		Radius: 1.0,
		Mat:    Dielectric{Refractance: 1.5}}
	diffuseSphere := Sphere{
		Center: mgl64.Vec3{-4, 1, 0},
		Radius: 1.0,
		Mat:    Lambertian{Albedo: mgl64.Vec3{0.4, 0.2, 0.1}}}
	metalSphere := Sphere{
		Center: mgl64.Vec3{4, 1, 0},
		Radius: 1.0,
		Mat:    Metallic{Albedo: mgl64.Vec3{0.4, 0.2, 0.1}, Fuzziness: 0.0}}

	scene = append(scene, dielectricSphere)
	scene = append(scene, diffuseSphere)
	scene = append(scene, metalSphere)

	return scene
}

// Render ray-traces the scene, outputting an image with the given dimensions.
func Render(width int, height int, samples int) image.Image {
	canvas := image.NewNRGBA(image.Rect(0, 0, width, height))

	from := mgl64.Vec3{13, 2, 3}
	at := mgl64.Vec3{0.0, 0.0, 0.0}
	fov := 20.0
	aperture := 0.1
	focusDistance := from.Sub(at).Len()
	cam := NewCamera(
		from,
		at,
		mgl64.Vec3{0, 1.0, 0.0},
		fov,
		float64(width)/float64(height),
		aperture,
		focusDistance)

	world := CreateRandomScene()

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
