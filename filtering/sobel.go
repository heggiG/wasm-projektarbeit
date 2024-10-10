package main

import (
	"image"
	"image/color"
	"math"
)

// sobel applies a sobel filter to a whole image
// The returned image has its size reduced by 2 in both dimensions.
func sobel(rgba image.Image) image.Image {
	grayImage := toGrayImage(rgba)
	convolved, low, high := sobelGray(grayImage)
	return toRGBAImage(convolved, low, high)
}

// toGrayImage turns an image gray
func toGrayImage(img image.Image) *image.Gray {
	grayImage := image.NewGray(img.Bounds())
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			grayImage.SetGray(x, y, grayColor)
		}
	}
	return grayImage
}

type kernel [9]int

// get small helper function to turn x-y coordinates into singular array access
func (k kernel) get(x, y int) int {
	return k[y*3+x]
}

// kernel for the vertical edge detection
var kernelX = kernel{
	1, 0, -1,
	2, 0, -2,
	1, 0, -1,
}

// kernel for the horizontal edge detection
var kernelY = kernel{
	1, 2, 1,
	0, 0, 0,
	-1, -2, -1,
}

// sobelGray applies the filtering to the whole image
func sobelGray(grayImage *image.Gray) (*image.Gray16, uint16, uint16) {
	width := grayImage.Bounds().Dx() - 2
	height := grayImage.Bounds().Dy() - 2
	convolved := image.NewGray16(image.Rect(0, 0, width, height))
	low, high := uint16(math.MaxUint16), uint16(0)
	for y := 1; y < grayImage.Bounds().Max.Y-1; y++ {
		for x := 1; x < grayImage.Bounds().Max.X-1; x++ {
			value := convolvePixel(grayImage, x, y)
			if low > value {
				low = value
			}
			if high < value {
				high = value
			}
			convolved.SetGray16(x-1, y-1, color.Gray16{Y: value})
		}
	}
	return convolved, low, high
}

// convolvePixel convolves a single pixel using both the x and y kernels
func convolvePixel(img *image.Gray, x, y int) uint16 {
	var valueX int
	var valueY int
	for ix, kx := x-1, 2; ix <= x+1; ix, kx = ix+1, kx-1 {
		for iy, ky := y-1, 2; iy <= y+1; iy, ky = iy+1, ky-1 {
			imgValue := int(img.GrayAt(ix, iy).Y)
			valueX += imgValue * kernelX.get(kx, ky)
			valueY += imgValue * kernelY.get(kx, ky)
		}
	}
	return uint16(math.Sqrt(float64(valueX*valueX + valueY*valueY)))
}

// toRGBAImage turns the image.Gray back into an image.RGBA
func toRGBAImage(img *image.Gray16, min uint16, max uint16) *image.RGBA {
	result := image.NewRGBA(img.Bounds())
	valueRange := float64(max - min)
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			value := img.Gray16At(x, y).Y
			outValue := uint8(float64(value-min) / valueRange * 255)
			result.Set(x, y, color.RGBA{R: outValue, G: outValue, B: outValue, A: 255})
		}
	}
	return result
}
