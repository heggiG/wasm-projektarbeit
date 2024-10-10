package main

import (
	"image"
	"image/color"
)

// blendColor shifts the pixel color towards a target color using a blend factor.
func blendColor(srcColor, targetColor color.Color, blendFactor float64) color.Color {
	r1, g1, b1, a1 := srcColor.RGBA()
	r2, g2, b2, a2 := targetColor.RGBA()

	r := uint8(float64(r1>>8)*(1-blendFactor) + float64(r2>>8)*blendFactor)
	g := uint8(float64(g1>>8)*(1-blendFactor) + float64(g2>>8)*blendFactor)
	b := uint8(float64(b1>>8)*(1-blendFactor) + float64(b2>>8)*blendFactor)
	a := uint8(float64(a1>>8)*(1-blendFactor) + float64(a2>>8)*blendFactor)

	return color.RGBA{R: r, G: g, B: b, A: a}
}

// colorShift applies blendColor to the whole image
func colorShift(img image.Image, targetColor color.Color, blendFactor float64) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Loop through each pixel
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			newColor := blendColor(originalColor, targetColor, blendFactor)
			newImg.Set(x, y, newColor)
		}
	}

	return newImg
}
