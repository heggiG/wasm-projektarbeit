package main

import (
	"florianheck/pong/pong"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"log"
)

// Game is the structure of the game state
type Game struct {
	state    pong.GameState
	cpuMode  bool
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

// Height of the game in pixels, not height of its window
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
	g.cpuMode = aiMode
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

// Update updates the game state based on player input and game logic.
func (g *Game) Update(screen *ebiten.Image) error {
	switch g.state {
	// if game is in start state check for key press and start the corresponding game or quit it
	case pong.StartState:
		if inpututil.IsKeyJustPressed(ebiten.KeyC) {
			g.state = pong.ControlsState
		} else if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.cpuMode = true
			g.state = pong.PlayState
		} else if inpututil.IsKeyJustPressed(ebiten.KeyV) {
			g.cpuMode = false
			g.state = pong.PlayState
		}

	// If game is in control state, pressing space will return to the main menu
	case pong.ControlsState:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.state = pong.StartState
		}

	// In play state, main game loop
	case pong.PlayState:

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

	err := g.Draw(screen)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// play is the game loop when in the playing state of the game
func (g *Game) play(screen *ebiten.Image) error {
	w, _ := screen.Size()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.state = pong.PauseState
		return nil
	}

	g.player1.Update(screen)
	if g.cpuMode {
		g.player2.AiUpdate(g.ball)
	} else {
		g.player2.Update(screen)
	}

	xV := g.ball.XVelocity
	g.ball.Update(g.player1, g.player2, screen)
	// rally count
	if xV*g.ball.XVelocity < 0 {
		// score up when ball touches human player's paddle
		if g.cpuMode && g.ball.X < float32(w/2) {
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
		if g.cpuMode {
			g.state = pong.GameOverState
			return nil
		}
		g.reset(screen, pong.InterState)
	} else if g.ball.X > float32(w) {
		g.player1.Score++
		if g.cpuMode {
			g.state = pong.GameOverState
			return nil
		}
		g.reset(screen, pong.InterState)
	}

	if g.player1.Score == g.maxScore || g.player2.Score == g.maxScore {
		g.state = pong.GameOverState
	}

	return nil
}

// Draw updates the game screen elements drawn
func (g *Game) Draw(screen *ebiten.Image) error {
	_ = screen.Fill(pong.BgColor)

	pong.DrawCaption(g.state, pong.ObjColor, screen)
	pong.DrawBigText(g.state, pong.ObjColor, screen)

	if g.state != pong.ControlsState {
		g.player1.Draw(screen, pong.ArcadeFont, false)
		g.player2.Draw(screen, pong.ArcadeFont, g.cpuMode)
		g.ball.Draw(screen)
	}

	// print the frames per second at the top left
	_ = ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))

	return nil
}

// Layout sets the screen layout
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func main() {
	// set the games window to fullscreen, makes it very easy to embed the game in an iframe
	ebiten.SetFullscreen(true)
	g := NewGame(true)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
