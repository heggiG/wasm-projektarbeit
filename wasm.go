package main

import (
	"bytes"
	"florianheck/wasm/pong"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image"
	"image/color"
	"image/png"
	"math/big"
	"strings"
	"syscall/js"
)

// Game is the structure of the game state
type Game struct {
	state    pong.GameState
	aiMode   bool
	ball     *pong.Ball
	player1  *pong.Paddle
	player2  *pong.Paddle
	rally    int
	level    int
	maxScore int
}

const (
	initBallVelocity = 5.0
	initPaddleSpeed  = 10.0
	speedUpdateCount = 6
	speedIncrement   = 0.5
)

const (
	windowWidth  = 800
	windowHeight = 600
)

// NewGame creates an initializes a new game
func NewGame(aiMode bool) *Game {
	g := &Game{}
	g.init(aiMode)
	return g
}

func (g *Game) init(aiMode bool) {
	g.state = pong.StartState
	g.aiMode = aiMode
	if aiMode {
		g.maxScore = 100
	} else {
		g.maxScore = 11
	}

	g.player1 = &pong.Paddle{
		Position: pong.Position{
			X: pong.InitPaddleShift,
			Y: float32(windowHeight / 2)},
		Score:  0,
		Speed:  initPaddleSpeed,
		Width:  pong.InitPaddleWidth,
		Height: pong.InitPaddleHeight,
		Color:  pong.ObjColor,
		Up:     ebiten.KeyW,
		Down:   ebiten.KeyS,
	}
	g.player2 = &pong.Paddle{
		Position: pong.Position{
			X: windowWidth - pong.InitPaddleShift - pong.InitPaddleWidth,
			Y: float32(windowHeight / 2)},
		Score:  0,
		Speed:  initPaddleSpeed,
		Width:  pong.InitPaddleWidth,
		Height: pong.InitPaddleHeight,
		Color:  pong.ObjColor,
		Up:     ebiten.KeyO,
		Down:   ebiten.KeyK,
	}
	g.ball = &pong.Ball{
		Position: pong.Position{
			X: float32(windowWidth / 2),
			Y: float32(windowHeight / 2)},
		Radius:    pong.InitBallRadius,
		Color:     pong.ObjColor,
		XVelocity: initBallVelocity,
		YVelocity: initBallVelocity,
	}
	g.level = 0
	g.ball.Img, _ = ebiten.NewImage(int(g.ball.Radius*2), int(g.ball.Radius*2), ebiten.FilterDefault)
	g.player1.Img, _ = ebiten.NewImage(g.player1.Width, g.player1.Height, ebiten.FilterDefault)
	g.player2.Img, _ = ebiten.NewImage(g.player2.Width, g.player2.Height, ebiten.FilterDefault)

	pong.InitFonts()
}

func (g *Game) reset(screen *ebiten.Image, state pong.GameState) {
	w, _ := screen.Size()
	g.state = state
	g.rally = 0
	g.level = 0
	if state == pong.StartState {
		g.player1.Score = 0
		g.player2.Score = 0
	}
	g.player1.Position = pong.Position{
		X: pong.InitPaddleShift, Y: pong.GetCenter(screen).Y}
	g.player2.Position = pong.Position{
		X: float32(w - pong.InitPaddleShift - pong.InitPaddleWidth), Y: pong.GetCenter(screen).Y}
	g.ball.Position = pong.GetCenter(screen)
	g.ball.XVelocity = initBallVelocity
	g.ball.YVelocity = initBallVelocity
}

// Update updates the game state
func (g *Game) Update(screen *ebiten.Image) error {
	switch g.state {
	case pong.StartState:
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			g.state = pong.ControlsState
		} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.aiMode = true
			g.state = pong.PlayState
		} else if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			g.aiMode = false
			g.state = pong.PlayState
		}

	case pong.ControlsState:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = pong.StartState
		}
	case pong.PlayState:
		w, _ := screen.Size()

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = pong.PauseState
			break
		}

		g.player1.Update(screen)
		if g.aiMode {
			g.player2.AiUpdate(g.ball)
		} else {
			g.player2.Update(screen)
		}

		xV := g.ball.XVelocity
		g.ball.Update(g.player1, g.player2, screen)
		// rally count
		if xV*g.ball.XVelocity < 0 {
			// score up when ball touches human player's paddle
			if g.aiMode && g.ball.X < float32(w/2) {
				g.player1.Score++
			}

			g.rally++

			// spice things up
			if (g.rally)%speedUpdateCount == 0 {
				g.level++
				g.ball.XVelocity += speedIncrement
				g.ball.YVelocity += speedIncrement
				g.player1.Speed += speedIncrement
				g.player2.Speed += speedIncrement
			}
		}

		if g.ball.X < 0 {
			g.player2.Score++
			if g.aiMode {
				g.state = pong.GameOverState
				break
			}
			g.reset(screen, pong.InterState)
		} else if g.ball.X > float32(w) {
			g.player1.Score++
			if g.aiMode {
				g.state = pong.GameOverState
				break
			}
			g.reset(screen, pong.InterState)
		}

		if g.player1.Score == g.maxScore || g.player2.Score == g.maxScore {
			g.state = pong.GameOverState
		}

	case pong.InterState, pong.PauseState:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = pong.PlayState
		} else if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.reset(screen, pong.StartState)
		}

	case pong.GameOverState:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.reset(screen, pong.StartState)
		}
	}

	g.Draw(screen)

	return nil
}

// Draw updates the game screen elements drawn
func (g *Game) Draw(screen *ebiten.Image) error {
	screen.Fill(pong.BgColor)

	pong.DrawCaption(g.state, pong.ObjColor, screen)
	pong.DrawBigText(g.state, pong.ObjColor, screen)

	if g.state != pong.ControlsState {
		g.player1.Draw(screen, pong.ArcadeFont, false)
		g.player2.Draw(screen, pong.ArcadeFont, g.aiMode)
		g.ball.Draw(screen)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))

	return nil
}

// Layout sets the screen layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

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
	png.Encode(&outputBuffer, resultImage)
	outputBytes := outputBuffer.Bytes()

	size := len(outputBytes)
	result := js.Global().Get("Uint8Array").New(size)
	js.CopyBytesToJS(result, outputBytes)

	return result
}

func main() {
	js.Global().Set("applySobel", js.FuncOf(ApplySobel))
	js.Global().Set("applyGaussean", js.FuncOf(ApplyGaussian))
	js.Global().Set("applyShift", js.FuncOf(ApplyShift))
	js.Global().Set("applyVignette", js.FuncOf(ApplyVignette))
	<-make(chan bool)
}
