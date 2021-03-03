package main

import (
	"bytes"
	"image/png"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten"
	"github.com/rpagliuca/go-physics/pkg/wave"
)

const SCREEN_WIDTH = 1000
const SCREEN_HEIGHT = 300

func main() {

	g := wave.Grid{}
	g[0][0] = 100
	g[0][1] = 95
	g[0][2] = 90
	g[0][3] = 85
	g[0][4] = 80
	g[1] = g[0]

	runGame(g)
}

func draw(grid wave.Grid) *ebiten.Image {
	dc := gg.NewContext(1000, 300)

	for i, u := range grid[2] {
		dc.DrawCircle(float64(i)*10, 300-2*u-100, 5)
		dc.SetRGB(0, 1.0, 0)
		dc.Fill()
	}

	writer := &bytes.Buffer{}
	dc.EncodePNG(writer)

	im, err := png.Decode(bytes.NewReader(writer.Bytes()))
	if err != nil {
		panic(err)
	}

	image, err := ebiten.NewImageFromImage(im, ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	return image
}

type Game struct {
	Grid wave.Grid
}

func (g *Game) Update(*ebiten.Image) error {
	g.Grid = wave.NextStep(g.Grid)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(draw(g.Grid), &ebiten.DrawImageOptions{})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return SCREEN_WIDTH, SCREEN_HEIGHT
}

func runGame(grid wave.Grid) {
	game := &Game{grid}
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(SCREEN_WIDTH, SCREEN_HEIGHT)
	ebiten.SetWindowTitle("Physics")
	// Call ebiten.RunGame to start your game loop.

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
