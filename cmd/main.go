package main

import (
	"flag"
	"image/png"
	"log"
	"os"

	"github.com/strangebroadcasts/weekendraytracer"
)

func main() {
	width := flag.Int("width", 200, "Width of image")
	height := flag.Int("height", 100, "Height of image")
	samples := flag.Int("samples", 8, "Samples for antialiasing")
	outputPath := flag.String("output", "out.png", "Output path")
	flag.Parse()
	image := weekendraytracer.Render(*width, *height, *samples)
	f, err := os.Create(*outputPath)
	if err != nil {
		log.Fatal("Error:", err)
	}
	err = png.Encode(f, image)
	if err != nil {
		f.Close()
		log.Fatal("Error:", err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal("Error:", err)
	}
}
