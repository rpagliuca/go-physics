package main

/* Source: https://raw.githubusercontent.com/cstegel/opengl-samples-golang/master/basic-shaders/main.go */

import (
	"log"
	"runtime"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	wave "github.com/rpagliuca/go-physics/pkg/wave-3d"

	"github.com/rpagliuca/go-gl-helpers/pkg/cam"
	"github.com/rpagliuca/go-gl-helpers/pkg/gfx"
	"github.com/rpagliuca/go-gl-helpers/pkg/win"
)

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to inifitialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window := win.NewWindow(1200, 800, "basic camera")

	// Initialize Glow (go function bindings)
	if err := gl.Init(); err != nil {
		panic(err)
	}

	err := programLoop(window)
	if err != nil {
		log.Fatal(err)
	}
}

/*
 * Creates the Vertex Array Object for a triangle.
 */
func createTriangleVAO(vertices []float32, indices []uint32) (uint32, func()) {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	/*
		var EBO uint32
		gl.GenBuffers(1, &EBO)
	*/

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STREAM_DRAW)

	// copy indices into element buffer
	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	//gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STREAM_DRAW)

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// normal
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	clear := func() {
		gl.DeleteVertexArrays(1, &VAO)
		gl.DeleteBuffers(1, &VBO)
	}

	return VAO, clear
}

func programLoop(window *win.Window) error {

	// the linked shader program determines how the data will be rendered
	vertShader, err := gfx.NewShaderFromFile("shaders/phong.vert", gl.VERTEX_SHADER)
	if err != nil {
		return err
	}

	fragShader, err := gfx.NewShaderFromFile("shaders/phong.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		return err
	}
	defer program.Delete()

	grid := wave.Grid{}

	/*
		// Center wave
		for i := 0; i < 6; i++ {
			for j := 0; j < 6; j++ {
				grid[0][wave.LEN/2-3+i][wave.LEN/2-3+j] = 10
			}
		}
	*/

	/*
		// Side wave
		for i := 0; i < wave.LEN; i++ {
			for j := 0; j < 5; j++ {
				grid[0][wave.LEN-1-j][i] = 2
			}
		}
	*/

	/*
		// Corner wave
		for i := 5; i < 15; i++ {
			for j := 5; j < 15; j++ {
				grid[0][wave.LEN-i][wave.LEN-j] = 5
			}
		}
	*/

	// 2 corner waves
	for i := 10; i < 30; i++ {
		for j := 10; j < 30; j++ {
			if (i+j)*(i+j) < 800 {
				grid[0][wave.LEN-i][wave.LEN-j] = 5
				grid[0][wave.LEN-i][j] = 5
			}
		}
	}

	grid[1] = grid[0]
	grid[2] = grid[0]

	program.Use()

	// ensure that triangles that are "behind" others do not draw over top of them
	//gl.Enable(gl.DEPTH_TEST)
	//gl.Enable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)
	//gl.DepthFunc(gl.LESS)
	//gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	camera := cam.NewFpsCamera(mgl32.Vec3{2.0, -2.0, 2.0}, mgl32.Vec3{0, 1, 0}, 45, 30, window.InputManager())

	for !window.ShouldClose() {

		//start := time.Now().UnixNano()
		window.StartFrame()
		camera.Update(window.SinceLastFrame())

		// perform rendering
		//gl.ClearColor(0.2, 0.5, 0.5, 1.0)
		//gl.ClearColor(135.0/255.0, 206.0/255.0, 250.0/255.0, 1.0)
		gl.ClearColor(0.1, 0.0, 0.3, 1.0)
		//glClearColor(255, 225, 175)
		//glClearColor(15, 68, 255)
		//glClearColor(163, 38, 82)
		//gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		drawScene(grid, window, camera, program)

		// end of draw loop

		grid = wave.NextStep(grid)
		//fmt.Println("window.shouldclose loop took milliseconds", (time.Now().UnixNano()-start)/1e6)
		//fmt.Println("current fps", math.Round(1.0/(float64(time.Now().UnixNano()-start)/1.0e9)))
	}

	return nil
}

func drawScene(waveGrid wave.Grid, window *win.Window, camera *cam.FpsCamera, program *gfx.Program) {

	//start := time.Now().UnixNano()

	grid := make([][]float64, len(waveGrid[2]))

	for i := 0; i < len(grid); i++ {
		grid[i] = make([]float64, len(grid))
	}

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid); j++ {
			grid[i][j] = waveGrid[2][i][j]
		}
	}

	refinedGrid := grid
	refinedGrid = refineGrid(grid)
	//refinedGrid = refineGrid(refinedGrid)
	//refinedGrid = refineGrid(refinedGrid)

	refinedPoint := func(i, j int) mgl32.Vec3 {
		x0 := i
		y0 := refinedGrid[i][j]
		z0 := j
		return mgl32.Vec3{float32(x0), float32(y0), float32(z0)}
	}
	//fmt.Println("half1 took milliseconds", (time.Now().UnixNano()-start)/1e6)

	//start = time.Now().UnixNano()
	vertices := make([]float32, len(refinedGrid)*len(refinedGrid)*2*3*6)

	var wg sync.WaitGroup

	for h := 0; h < len(refinedGrid)-1; h++ {
		wg.Add(1)
		go func(i int) {
			for j := 0; j < len(refinedGrid)-1; j++ {
				pos := 36 * (i*len(refinedGrid) + j)
				vertices = addVertices(vertices, pos, refinedPoint(i, j), refinedPoint(i+1, j+1), refinedPoint(i+1, j))
				vertices = addVertices(vertices, pos+18, refinedPoint(i, j), refinedPoint(i, j+1), refinedPoint(i+1, j+1))
			}
			wg.Done()
		}(h)
	}
	wg.Wait()

	//fmt.Println("addVertices took milliseconds", (time.Now().UnixNano()-start)/1e6)

	//start = time.Now().UnixNano()

	//VAO := createTriangleVAO(vertices, indices)
	VAO, clearVAOFunc := createTriangleVAO(vertices, nil)

	// creates perspective
	fov := float32(60.0)
	projectTransform := mgl32.Perspective(mgl32.DegToRad(fov),
		float32(window.Width())/float32(window.Height()),
		0.01,
		100.0)

	camTransform := camera.GetTransform()
	worldTransform := mgl32.Scale3D(0.05, 0.2, 0.05)
	gl.UniformMatrix4fv(program.GetUniformLocation("project"), 1, false, &projectTransform[0])
	gl.UniformMatrix4fv(program.GetUniformLocation("view"), 1, false, &camTransform[0])
	gl.UniformMatrix4fv(program.GetUniformLocation("world"), 1, false, &worldTransform[0])

	gl.BindVertexArray(VAO)
	//gl.DrawArrays(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))
	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))
	//gl.DrawArrays(gl.TRIANGLES, 0, 6)
	gl.BindVertexArray(0)
	clearVAOFunc()
	//fmt.Println("half2 took milliseconds", (time.Now().UnixNano()-start)/1e6)
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action,
	mods glfw.ModifierKey) {

	// When a user presses the escape key, we set the WindowShouldClose property to true,
	// which closes the application
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func interpolation(x1, y1, x2, y2, x3, y3 float32) func(float32) float32 {
	return func(x float32) float32 {
		t1 := y1 * (x - x2) * (x - x3) / ((x1 - x2) * (x1 - x3))
		t2 := y2 * (x - x1) * (x - x3) / ((x2 - x1) * (x2 - x3))
		t3 := y3 * (x - x1) * (x - x2) / ((x3 - x1) * (x3 - x2))
		return t1 + t2 + t3
	}
}

// p_0 and p_1 clockwise
func normal(p_center, p_0, p_1 mgl32.Vec3) (x, y, z float32) {
	return (p_0.Sub(p_center)).Cross(p_1.Sub(p_center)).Normalize().Mul(-1.0).Elem()
}

func addVertices(vertices []float32, pos int, p0, p1, p2 mgl32.Vec3) []float32 {

	x, y, z := normal(p0, p1, p2)

	vertices[pos+0] = p1.X()
	vertices[pos+1] = p1.Y()
	vertices[pos+2] = p1.Z()
	vertices[pos+3] = x
	vertices[pos+4] = y
	vertices[pos+5] = z
	vertices[pos+6] = p2.X()
	vertices[pos+7] = p2.Y()
	vertices[pos+8] = p2.Z()
	vertices[pos+9] = x
	vertices[pos+10] = y
	vertices[pos+11] = z
	vertices[pos+12] = p0.X()
	vertices[pos+13] = p0.Y()
	vertices[pos+14] = p0.Z()
	vertices[pos+15] = x
	vertices[pos+16] = y
	vertices[pos+17] = z

	return vertices
}

func q(x float64) float64 {
	return x * x
}

func quadratic(x0, y0, x1, y1, x2, y2 float64) func(float64) float64 {
	b := ((y2-y0)*(q(x1)-q(x0)) - (y1-y0)*(q(x2)*q(x0))) / ((x2-x0)*(q(x1)-q(x0)) + (x1-x0)*(q(x2)-q(x0)))
	a := (y2 - y0 - b*(x2-x0)) / (q(x2) - q(x0))
	c := y2 - a*q(x2) - b*x2
	return func(x float64) float64 {
		return a*q(x) + b*x + c
	}
}

func refineGrid(grid [][]float64) [][]float64 {

	// Initialize 2D slice
	refinedGrid := make([][]float64, 2*len(grid))
	for i := 0; i < len(refinedGrid); i++ {
		refinedGrid[i] = make([]float64, 2*len(grid))
	}

	for i := 0; i < len(grid)-2; i++ {
		for j := 0; j < len(grid)-2; j++ {

			p_c := grid[i][j]

			p_t := grid[i][j+1]
			p_t2 := grid[i][j+2]
			p_tr := grid[i+1][j+1]
			p_tr2 := grid[i+2][j+2]
			p_r := grid[i+1][j]
			p_r2 := grid[i+2][j]

			f_ct := quadratic(0, p_c, 1, p_t, 2, p_t2)
			f_cr := quadratic(0, p_c, 1, p_r, 2, p_r2)
			f_ctr := quadratic(0, p_c, 1, p_tr, 2, p_tr2)

			p_cr := f_cr(0.5)
			p_ct := f_ct(0.5)
			p_ctr := f_ctr(0.5)

			refinedGrid[2*i][2*j] = p_c
			refinedGrid[2*i+1][2*j] = p_cr
			refinedGrid[2*i][2*j+1] = p_ct
			refinedGrid[2*i+1][2*j+1] = p_ctr
		}
	}
	return refinedGrid
}

func glClearColor(r, g, b int) {
	gl.ClearColor(
		float32(r)/255.0,
		float32(g)/255.0,
		float32(b)/255.0,
		1.0,
	)
}
