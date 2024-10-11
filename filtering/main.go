package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"strings"
	"syscall/js"
)

// GetColorFromString is a helper method to get a color.RGBA value from a string in the #FFFFFF style
func GetColorFromString(hex string) color.RGBA {
	n := new(big.Int)
	n.SetString(strings.TrimPrefix(hex, "#"), 16)
	val := n.Int64()
	c := color.RGBA{R: uint8((val & 0xFF0000) >> 16), G: uint8((val & 0xFF00) >> 8), B: uint8(val & 0xFF), A: 255}
	return c
}

// ApplySobel calls ApplyImageOperator to apply a sobel filter to the image
func ApplySobel(this js.Value, args []js.Value) interface{} {
	return ApplyImageOperator(this, args, "sobel")
}

// ApplyGaussian calls ApplyImageOperator to apply a gaussian blur to the image
func ApplyGaussian(this js.Value, args []js.Value) interface{} {
	return ApplyImageOperator(this, args, "gaussian")
}

// ApplyShift calls ApplyImageOperator to apply a color shift to the image
func ApplyShift(this js.Value, args []js.Value) interface{} {
	return ApplyImageOperator(this, args, "shift")
}

// ApplyVignette calls ApplyImageOperator to apply a vignette to the image
func ApplyVignette(this js.Value, args []js.Value) interface{} {
	return ApplyImageOperator(this, args, "vignette")
}

// ApplyImageOperator creates the image from a buffer and applies the corresponding operation to the image
func ApplyImageOperator(this js.Value, args []js.Value, operation string) interface{} {
	inputBuffer := make([]byte, args[0].Get("byteLength").Int())
	js.CopyBytesToGo(inputBuffer, args[0])
	img, _, _ := image.Decode(bytes.NewReader(inputBuffer))

	var resultImage image.Image

	switch operation {
	case "sobel":
		resultImage = sobel(img)
		break

	case "gaussian":
		resultImage = applyGaussianBlur(img)
		break

	case "shift":
		resultImage = colorShift(img, GetColorFromString(args[1].String()), args[2].Float())
		break

	case "vignette":
		resultImage = addVignette(img, args[1].Float(), args[2].Float(), args[3].Float(), args[4].Float())

	default:
		panic("No valid operation given to execute")
	}

	if resultImage == nil {
		return nil
	}

	var outputBuffer bytes.Buffer
	err := png.Encode(&outputBuffer, resultImage)
	if err != nil {
		panic(err)
	}
	outputBytes := outputBuffer.Bytes()

	size := len(outputBytes)
	result := js.Global().Get("Uint8Array").New(size)
	js.CopyBytesToJS(result, outputBytes)

	return result
}

// The main method calls js to globally set the functions written in go as javascript methods so they can be called from the dom.
func main() {
	js.Global().Set("applySobel", js.FuncOf(ApplySobel))
	js.Global().Set("applyGaussian", js.FuncOf(ApplyGaussian))
	js.Global().Set("applyShift", js.FuncOf(ApplyShift))
	js.Global().Set("applyVignette", js.FuncOf(ApplyVignette))
	<-make(chan bool)
}
