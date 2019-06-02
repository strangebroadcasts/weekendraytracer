package weekendraytracer

import (
	"image"
	"image/color"

	"github.com/go-gl/mathgl/mgl64"
)

// Render ray-traces the scene, outputting an image with the given dimensions.
func Render(width int, height int) image.Image {
	canvas := image.NewNRGBA(image.Rect(0, 0, width, height))
	// Since the image library's axes increase right and down,
	// we don't need to flip the vertical axis to write rows
	// top to bottom, as RTiaW does:
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			col := mgl64.Vec3{float64(i) / float64(width), float64(j) / float64(height), 0.2}
			ir, ig, ib := uint8(255.99*col.X()), uint8(255.99*col.Y()), uint8(255.99*col.Z())
			canvas.Set(i, j, color.NRGBA{R: ir, G: ig, B: ib, A: 255})
		}
	}
	return canvas
}
