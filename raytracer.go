package weekendraytracer

import (
	"image"
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// Get the pixel color for this ray.
// (Called "color" in RTiaW, which conflicts with the color package)
func getColor(r Ray, world HittableList) mgl64.Vec3 {
	hits := world.Hit(r, 0.0, math.MaxFloat64)
	if len(hits) > 0 {
		// HittableList makes sure the first (and only) intersection
		// is the closest one:
		N := hits[0].Normal
		return N.Add(mgl64.Vec3{1.0, 1.0, 1.0}).Mul(0.5)
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
func Render(width int, height int) image.Image {
	canvas := image.NewNRGBA(image.Rect(0, 0, width, height))

	lowerLeftCorner := mgl64.Vec3{-2.0, -1.0, -1.0}
	horizontal := mgl64.Vec3{4.0, 0.0, 0.0}
	vertical := mgl64.Vec3{0.0, 2.0, 0.0}
	origin := mgl64.Vec3{0.0, 0.0, 0.0}

	world := make(HittableList, 2)
	world[0] = Sphere{Center: mgl64.Vec3{0.0, 0.0, -1.0}, Radius: 0.5}
	world[1] = Sphere{Center: mgl64.Vec3{0.0, -100.5, -1.0}, Radius: 100}

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			// Note that we flip the vertical axis here.
			u, v := float64(i)/float64(width), float64(height-j)/float64(height)
			rayDirection := lowerLeftCorner.Add(horizontal.Mul(u).Add(vertical.Mul(v)))
			r := Ray{A: origin, B: rayDirection}
			col := getColor(r, world)
			ir, ig, ib := uint8(255.99*col.X()), uint8(255.99*col.Y()), uint8(255.99*col.Z())
			canvas.Set(i, j, color.NRGBA{R: ir, G: ig, B: ib, A: 255})
		}
	}
	return canvas
}
