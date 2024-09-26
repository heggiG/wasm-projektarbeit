package main

import (
	"bytes"
	"image"
	"image/png"
	"syscall/js"
)

func applySobel(this js.Value, args []js.Value) interface{} {
	return applyImageOperator(this, args, "sobel")
}

func applyGaussean(this js.Value, args []js.Value) interface{} {
	return applyImageOperator(this, args, "gaussean")
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

	default:
		panic("No valid operation given to execute")
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
	<-make(chan bool)
}
