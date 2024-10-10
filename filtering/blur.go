package main

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
)

// Kernel used in the gaussian blur
var gaussianKernel = [11][11]float64{
	{0.0000, 0.0000, 0.0000, 0.0001, 0.0003, 0.0003, 0.0003, 0.0001, 0.0000, 0.0000, 0.0000},
	{0.0000, 0.0001, 0.0003, 0.0009, 0.0018, 0.0022, 0.0018, 0.0009, 0.0003, 0.0001, 0.0000},
	{0.0000, 0.0003, 0.0014, 0.0042, 0.0080, 0.0099, 0.0080, 0.0042, 0.0014, 0.0003, 0.0000},
	{0.0001, 0.0009, 0.0042, 0.0123, 0.0234, 0.0290, 0.0234, 0.0123, 0.0042, 0.0009, 0.0001},
	{0.0003, 0.0018, 0.0080, 0.0234, 0.0445, 0.0551, 0.0445, 0.0234, 0.0080, 0.0018, 0.0003},
	{0.0003, 0.0022, 0.0099, 0.0290, 0.0551, 0.0682, 0.0551, 0.0290, 0.0099, 0.0022, 0.0003},
	{0.0003, 0.0018, 0.0080, 0.0234, 0.0445, 0.0551, 0.0445, 0.0234, 0.0080, 0.0018, 0.0003},
	{0.0001, 0.0009, 0.0042, 0.0123, 0.0234, 0.0290, 0.0234, 0.0123, 0.0042, 0.0009, 0.0001},
	{0.0000, 0.0003, 0.0014, 0.0042, 0.0080, 0.0099, 0.0080, 0.0042, 0.0014, 0.0003, 0.0000},
	{0.0000, 0.0001, 0.0003, 0.0009, 0.0018, 0.0022, 0.0018, 0.0009, 0.0003, 0.0001, 0.0000},
	{0.0000, 0.0000, 0.0000, 0.0001, 0.0003, 0.0003, 0.0003, 0.0001, 0.0000, 0.0000, 0.0000},
}

// applyGaussianBlur applies a 7x7 Gaussian blur to the given image
func applyGaussianBlur(img image.Image) image.Image {
	bounds := img.Bounds()
	blurred := image.NewRGBA(bounds)
	draw.Draw(blurred, bounds, img, bounds.Min, draw.Src)

	kernelSize := len(gaussianKernel)
	offset := kernelSize / 2

	for x := bounds.Min.X + offset; x < bounds.Max.X-offset; x++ {
		for y := bounds.Min.Y + offset; y < bounds.Max.Y-offset; y++ {
			var r, g, b float64

			// Convolve the pixel with the kernel
			for i := 0; i < kernelSize; i++ {
				for j := 0; j < kernelSize; j++ {
					px := img.At(x-offset+i, y-offset+j)
					rgba := color.RGBAModel.Convert(px).(color.RGBA)
					weight := gaussianKernel[i][j]

					r += float64(rgba.R) * weight
					g += float64(rgba.G) * weight
					b += float64(rgba.B) * weight
				}
			}

			// Set the new pixel in the blurred image
			blurred.Set(x, y, color.RGBA{
				R: uint8(clamp(r, 0, 255)),
				G: uint8(clamp(g, 0, 255)),
				B: uint8(clamp(b, 0, 255)),
				A: 255, // Set alpha to fully opaque
			})
		}
	}

	return blurred
}

// clamp ensures a value is within a given range
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
