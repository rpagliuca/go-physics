package main

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/rpagliuca/go-physics/pkg/wave"
)

func main() {

	g := wave.Grid{}
	g[0][0] = 100
	g[0][1] = 90
	g[0][2] = 80
	g[0][3] = 70
	g[0][4] = 60
	g[0][5] = 50
	g[1] = g[0]

	for i := 0; i < 2000; i++ {
		g = wave.NextStep(g)
		draw(i, g)
	}

}

func draw(step int, grid wave.Grid) {
	dc := gg.NewContext(1000, 300)

	for i, u := range grid[2] {
		dc.DrawCircle(float64(i)*10, 300-2*u-150, 5)
		dc.SetRGB(0, 1.0, 0)
		dc.Fill()
	}

	dc.SavePNG(fmt.Sprintf("out/out_%02d.png", step))
}
