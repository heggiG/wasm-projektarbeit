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

func getColorFromString(hex string) color.RGBA {
	n := new(big.Int)
	n.SetString(strings.TrimPrefix(hex, "#"), 16)
	val := n.Int64()
	c := color.RGBA{R: uint8((val & 0xFF0000) >> 16), G: uint8((val & 0xFF00) >> 8), B: uint8(val & 0xFF), A: 255}
	return c
}

func applySobel(this js.Value, args []js.Value) interface{} {
	return applyImageOperator(this, args, "sobel")
}

func applyGaussean(this js.Value, args []js.Value) interface{} {
	return applyImageOperator(this, args, "gaussean")
}

func applyShift(this js.Value, args []js.Value) interface{} {
	return applyImageOperator(this, args, "shift")
}

func applyVignette(this js.Value, args []js.Value) interface{} {
	return applyImageOperator(this, args, "vignette")
}

func applyImageOperator(this js.Value, args []js.Value, operation string) interface{} {
	inputBuffer := make([]byte, args[0].Get("byteLength").Int())
	js.CopyBytesToGo(inputBuffer, args[0])
	img, _, _ := image.Decode(bytes.NewReader(inputBuffer))

	var resultImage image.Image

	switch operation {
	case "sobel":
		resultImage = sobel(img)
		break

	case "gaussean":
		resultImage = applyGaussianBlur(img)
		break

	case "shift":
		resultImage = colorShift(img, getColorFromString(args[1].String()), args[2].Float())
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
	png.Encode(&outputBuffer, resultImage)
	outputBytes := outputBuffer.Bytes()

	size := len(outputBytes)
	result := js.Global().Get("Uint8Array").New(size)
	js.CopyBytesToJS(result, outputBytes)

	return result
}

func main() {
	js.Global().Set("applySobel", js.FuncOf(applySobel))
	js.Global().Set("applyGaussean", js.FuncOf(applyGaussean))
	js.Global().Set("applyShift", js.FuncOf(applyShift))
	js.Global().Set("applyVignette", js.FuncOf(applyVignette))
	<-make(chan bool)
}
