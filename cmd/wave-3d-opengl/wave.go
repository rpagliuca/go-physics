package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	wave "github.com/rpagliuca/go-physics/pkg/wave-3d"
)

const SCREEN_WIDTH = 1200
const SCREEN_HEIGHT = 600

const width = SCREEN_WIDTH
const height = SCREEN_HEIGHT

var (
	rotationX float32
	rotationY float32
	camz      = float32(-60.0)
	camx      = float32(2.0)
	camy      = float32(5.0)
	drz       = float32(30.0)
	drx       = float32(0.0)
	dry       = float32(0.0)
)

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(width, height, "go-physics", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	grid := wave.Grid{}

	grid[0][0][0] = 20
	grid[0][1][1] = 15
	grid[0][2][2] = 10
	grid[0][3][3] = 5

	grid[1] = grid[0]
	grid[2] = grid[0]

	for !window.ShouldClose() {

		setupScene()

		drawScene(grid)
		window.SwapBuffers()
		glfw.PollEvents()
		grid = wave.NextStep(grid)
	}
}

func setupScene() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)

	gl.ClearColor(1, 1, 1, 1.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	ambient := []float32{0.0, 0.0, 0.1, 0.0}
	diffuse := []float32{0.0, 0.0, 50.0, 5.0}
	lightPosition := []float32{-1, -1, 0.5, 1}
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	gl.Enable(gl.LIGHT0)

	gl.MatrixMode(gl.PROJECTION)

	gl.LoadIdentity()
	//gl.UniformMatrix4fv(uniform, 1, false, &eye[0])
	//gl.LoadMatrixf(&eye[0])
	gl.Frustum(-1, 1, -0.7, 0.05, 4.0, 500.0)

	eye := mgl32.LookAt(50.0, 50.0, 20.0, 0, 0, 0, 0, 0, 1.0)
	gl.MultMatrixf(&eye[0])

	//gl.Translatef(0, 0, -40)

	/*
		drz += 0.10
		dry += 0.5
		drx += 0.25
		gl.Rotatef(drz, 1, 0, 0)
		gl.Rotatef(dry, 0, 1, 0)
		gl.Rotatef(drx, 0, 0, 1)
	*/

	/*
		gl.Rotatef(-90, 0, 0, 1)
		gl.Translatef(-20, -20, 0)
		/*
			gl.Rotatef(20, 1, 0, 0)
			gl.Rotatef(-20, 0, 1, 0)
			gl.Rotatef(10, 0, 0, 1)
			gl.Rotatef(20, 1, 0, 0)
	*/
	//gl.Rotatef(90, 0, 1, 1)
	//gl.Translatef(10, 0, 0)
	/*
		drz -= 0.02
		gl.Rotatef(drz, 1, 0, 0)

		drx -= 0.02
		gl.Rotatef(drx, 0, 1, 0)
	*/

	/*
		dry -= 0.02
		gl.Rotatef(dry, 0, 0, 1)
	*/

	//fmt.Println(drz, drx, dry)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func destroyScene() {
}

func drawScene(grid wave.Grid) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	//gl.Translatef(0, 0, -3.0)
	//gl.Rotatef(-40, 1, 0, 0)

	for i := 0; i < wave.LEN; i++ {
		for j := 0; j < wave.LEN; j++ {
			drawPlane(float32(i)/5.0, float32(j)/5.0, float32(grid[2][i][j]*2.0))
		}
	}

}

func drawPlane(x, y, z float32) {

	/*
		if z < 0 {
			z = 0
		}
	*/

	gl.Begin(gl.QUADS)

	h := float32(0.1)

	gl.Normal3f(0, 0, 1)

	color1 := []float32{0.0, 0.0, 0.0}
	//color2 := []float32{0.0, 1.0, 0.0}

	//gl.Materialfv(gl.FRONT, gl.AMBIENT, &color2[0])
	//gl.Materialfv(gl.FRONT, gl.DIFFUSE, &color2[0])
	shin := float32(30.0)
	gl.Materialfv(gl.FRONT, gl.SPECULAR, &color1[0])
	gl.Materialfv(gl.FRONT, gl.SHININESS, &shin)

	gl.Vertex3f(x-h, y-h, z)
	gl.Vertex3f(x+h, y-h, z)
	gl.Vertex3f(x+h, y+h, z)
	gl.Vertex3f(x-h, y+h, z)

	gl.Vertex3f(x-h, y-h, -1)
	gl.Vertex3f(x+h, y-h, -1)
	gl.Vertex3f(x+h, y+h, -1)
	gl.Vertex3f(x-h, y+h, -1)

	gl.Vertex3f(x-h, y-h, -1)
	gl.Vertex3f(x+h, y-h, -1)
	gl.Vertex3f(x+h, y-h, z)
	gl.Vertex3f(x-h, y-h, z)

	gl.Vertex3f(x-h, y-h, -1)
	gl.Vertex3f(x-h, y+h, -1)
	gl.Vertex3f(x-h, y+h, z)
	gl.Vertex3f(x-h, y-h, z)

	gl.Vertex3f(x+h, y-h, -1)
	gl.Vertex3f(x+h, y+h, -1)
	gl.Vertex3f(x+h, y+h, z)
	gl.Vertex3f(x+h, y-h, z)

	gl.Vertex3f(x+h, y+h, -1)
	gl.Vertex3f(x-h, y+h, -1)
	gl.Vertex3f(x-h, y+h, z)
	gl.Vertex3f(x+h, y+h, z)

	gl.End()
}
