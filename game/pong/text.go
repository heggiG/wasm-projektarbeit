package pong

import (
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"log"
)

const (
	fontSize      = 30
	smallFontSize = int(fontSize / 2)
)

var (
	ArcadeFont      font.Face
	SmallArcadeFont font.Face
)

// InitFonts initiales the font used in the game
func InitFonts() {
	tt, err := truetype.Parse(fonts.ArcadeN_ttf)
	if err != nil {
		log.Fatal(err)
	}
	var dpi float64 = 72
	ArcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(fontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	SmallArcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    float64(smallFontSize),
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

// DrawCaption draws the smaller text at the bottom of the screen
func DrawCaption(state GameState, color color.Color, screen *ebiten.Image) {
	w, h := screen.Size()
	var msg []string
	switch state {
	case PlayState, InterState, PauseState:
		msg = append(msg, "Press SPACE key to take a break (not too long though)")
	case ControlsState:
		msg = append(msg, "Press SPACE to go back to main menu")
	default:
		// noop
	}
	for i, l := range msg {
		x := (w - len(l)*smallFontSize) / 2
		text.Draw(screen, l, SmallArcadeFont, x, h-4+(i-2)*smallFontSize, color)
	}
}

// DrawBigText shows large text on the, depending on game state
func DrawBigText(state GameState, color color.Color, screen *ebiten.Image) {
	w, _ := screen.Size()
	var texts []string
	switch state {
	case StartState:
		texts = []string{
			"",
			"PONG",
			"",
			"C -> CONTROLS",
			"V -> VS GAME",
			"A -> CPU GAME",
		}
	case ControlsState:
		texts = []string{
			"",
			"PLAYER 1:",
			"W -> UP",
			"S -> DOWN",
			"",
			"PLAYER 2:",
			"O -> UP",
			"K -> DOWN",
		}
	case PauseState:
		texts = []string{
			"",
			"PAUSED",
			"",
			"SPACE -> RESUME",
			"R     -> RESET",
		}
	case InterState:
		texts = []string{
			"",
			"",
			"SPACE -> RESUME",
			"R     -> RESET",
		}
	case GameOverState:
		texts = []string{
			"",
			"GAME OVER!",
			"SPACE -> RESET",
		}
	default:
		// noop
	}
	for i, l := range texts {
		x := (w - len(l)*fontSize) / 2
		text.Draw(screen, l, ArcadeFont, x, (i+4)*fontSize, color)
	}
}
