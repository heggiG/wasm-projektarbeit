package main

import (
	"image"
	"image/color"
	"math"
)

func addVignette(img image.Image, radius float64, centerX, centerY, strength float64) image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	if int(centerX) > width || int(centerY) > height {
		return nil
	}

	// Create a new image to apply the vignette effect
	vignettedImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < height; y++ {
		for x := bounds.Min.X; x < width; x++ {
			// Calculate the distance from the pixel to the vignette center
			distX := float64(x) - centerX
			distY := float64(y) - centerY
			distance := math.Sqrt(distX*distX + distY*distY)

			// Calculate the vignette factor based on distance
			// Pixels within the radius are not darkened; beyond the radius they are darkened based on distance
			factor := 1.0
			if distance > radius {
				factor = 1.0 - math.Min(1.0, (distance-radius)/(radius*strength))
			}

			// Get the original pixel color
			origColor := img.At(x, y)
			r, g, b, a := origColor.RGBA()

			// Apply the vignette factor to the pixel color
			newR := uint8(float64(r>>8) * factor)
			newG := uint8(float64(g>>8) * factor)
			newB := uint8(float64(b>>8) * factor)

			// Set the new color in the resulting image
			vignettedImg.Set(x, y, color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)})
		}
	}

	return vignettedImg
}
