package pong

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

// Ball is a game ball
type Ball struct {
	Position
	Radius    float32
	XVelocity float32
	YVelocity float32
	Color     color.Color
	Img       *ebiten.Image
}

const (
	InitBallRadius = 10.0
)

func setBallPixels(c color.Color, ballImg *ebiten.Image) {
	err := ballImg.Fill(c)
	if err != nil {
		panic(err)
	}
}

func (ball *Ball) Update(leftPaddle *Paddle, rightPaddle *Paddle, screen *ebiten.Image) {
	_, h := screen.Size()
	ball.X += ball.XVelocity
	ball.Y += ball.YVelocity

	// bounce off edges when getting to top/bottom
	if ball.Y-ball.Radius > float32(h) {
		ball.YVelocity = -ball.YVelocity
		ball.Y = float32(h) - ball.Radius
	} else if ball.Y+ball.Radius < 0 {
		ball.YVelocity = -ball.YVelocity
		ball.Y = ball.Radius
	}

	// bounce off paddles
	if ball.X-ball.Radius < leftPaddle.X+float32(leftPaddle.Width/2) &&
		ball.Y > leftPaddle.Y-float32(leftPaddle.Height/2) &&
		ball.Y < leftPaddle.Y+float32(leftPaddle.Height/2) {
		ball.XVelocity = -ball.XVelocity
		ball.X = leftPaddle.X + float32(leftPaddle.Width/2) + ball.Radius
	} else if ball.X+ball.Radius > rightPaddle.X-float32(rightPaddle.Width/2) &&
		ball.Y > rightPaddle.Y-float32(rightPaddle.Height/2) &&
		ball.Y < rightPaddle.Y+float32(rightPaddle.Height/2) {
		ball.XVelocity = -ball.XVelocity
		ball.X = rightPaddle.X - float32(rightPaddle.Width/2) - ball.Radius
	}
}

func (ball *Ball) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(ball.X), float64(ball.Y))
	setBallPixels(ball.Color, ball.Img)
	err := screen.DrawImage(ball.Img, opts)
	if err != nil {
		panic(err)
	}
}
