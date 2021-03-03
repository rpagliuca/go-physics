package main

import (
	"io/ioutil"
	"log"

	"github.com/fogleman/ln/ln"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	wave "github.com/rpagliuca/go-physics/pkg/wave-3d"
)

const SCREEN_WIDTH = 500
const SCREEN_HEIGHT = 500

func main() {

	g0 := wave.Grid{}

	g0[0][0][0] = 10
	g0[0][1][1] = 9
	g0[0][2][2] = 8
	g0[0][3][3] = 7

	g0[1] = g0[0]
	g0[2] = g0[0]

	runGame(g0)
}

func draw(grid wave.Grid) *ebiten.Image {
	// create a scene and add a single cube
	scene := ln.Scene{}

	for i := 0; i < wave.LEN; i++ {
		for j := 0; j < wave.LEN; j++ {
			size := 1.0
			x := float64(i)
			y := float64(j)
			scene.Add(ln.NewCube(ln.Vector{x, y, grid[2][i][j]}, ln.Vector{x + size, y + size, grid[2][i][j] + size}))
			//scene.Add(ln.NewSphere(ln.Vector{x, y, grid[2][i][j]}, size))
			//scene.Add(ln.NewCube(ln.Vector{x, y, grid[2][i][j]}, ln.Vector{x + size, y + size, grid[2][i][j] + size}))
		}
	}

	// define camera parameters
	eye := ln.Vector{30, 30, 10}                       // camera position
	center := ln.Vector{wave.LEN / 2, wave.LEN / 2, 0} // camera looks at
	up := ln.Vector{0, 0, 1}                           // up direction

	// define rendering parameters
	width := float64(SCREEN_WIDTH)   // rendered width
	height := float64(SCREEN_HEIGHT) // rendered height
	fovy := 50.0                     // vertical field of view, degrees
	znear := 0.1                     // near z plane
	zfar := 50.0                     // far z plane
	step := 0.01                     // how finely to chop the paths for visibility testing

	// compute 2D paths that depict the 3D scene
	paths := scene.Render(eye, center, up, width, height, fovy, znear, zfar, step)

	// render the paths in an image
	file, err := ioutil.TempFile("", "")

	if err != nil {
		panic(err)
	}

	filepath := file.Name()

	paths.WriteToPNG(filepath, width, height)

	image, _, err := ebitenutil.NewImageFromFile(filepath, ebiten.FilterDefault)
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
	ebiten.SetMaxTPS(20)
	// Call ebiten.RunGame to start your game loop.

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
