package main

import (
	"image"
	"image/color"
	"math"
)

// distance returns the distance between two points (x1, y1) and (x2, y2)
func distance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))
}

// Vignette applies a vignette effect to an image based on the provided parameters
func Vignette(img image.Image, centerX, centerY, radius int, strength float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Create a new image to apply the vignette effect
	newImg := image.NewRGBA(bounds)

	// Calculate the maximum possible distance from the center
	maxDist := distance(0, 0, width, height)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the distance of the current pixel from the center point
			dist := distance(x, y, centerX, centerY)

			// Calculate vignette factor (based on distance and radius)
			vignetteFactor := 1.0
			if dist > float64(radius) {
				// The farther a pixel is from the center, the stronger the vignette effect
				vignetteFactor = 1 - strength*(math.Min(1.0, (dist-float64(radius))/(maxDist-float64(radius))))
			}

			// Get the original color of the pixel
			r, g, b, a := img.At(x, y).RGBA()

			// Apply the vignette factor (darken the pixel)
			newR := uint8(float64(r>>8) * vignetteFactor)
			newG := uint8(float64(g>>8) * vignetteFactor)
			newB := uint8(float64(b>>8) * vignetteFactor)

			// Set the new pixel color
			newImg.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return newImg
}
