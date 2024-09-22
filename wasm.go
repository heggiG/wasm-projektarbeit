package main

import (
	"bytes"
	"image"
	"image/png"
	"syscall/js"
)

func applySobelOperator(this js.Value, args []js.Value) any {
    inputBuffer := make([]byte, args[0].Get("byteLength").Int())
    js.CopyBytesToGo(inputBuffer, args[0])
    img, _, _ := image.Decode(bytes.NewReader(inputBuffer))
    
    resultImage := sobel(img)
    
    var outputBuffer bytes.Buffer
    png.Encode(&outputBuffer, resultImage)
    outputBytes := outputBuffer.Bytes()
    
    size := len(outputBytes)
    result := js.Global().Get("Uint8Array").New(size)
    js.CopyBytesToJS(result, outputBytes)
    
    return result
}

func main() {
    js.Global().Set("applySobel", js.FuncOf(applySobelOperator))
	<-make(chan bool);
}

