package main

/* Source: https://raw.githubusercontent.com/cstegel/opengl-samples-golang/master/basic-shaders/main.go */

import (
	"log"
	"runtime"

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

	// Corner wave
	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			grid[0][wave.LEN-i][wave.LEN-j] = 5
		}
	}

	/*
		// 2 corner waves
		for i := 10; i < 30; i++ {
			for j := 10; j < 30; j++ {
				if (i+j)*(i+j) < 800 {
					grid[0][wave.LEN-i][wave.LEN-j] = 5
					grid[0][wave.LEN-i][j] = 5
				}
			}
		}
	*/

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
		window.StartFrame()
		camera.Update(window.SinceLastFrame())

		// perform rendering
		//gl.ClearColor(0.2, 0.5, 0.5, 1.0)
		//gl.ClearColor(135.0/255.0, 206.0/255.0, 250.0/255.0, 1.0)
		gl.ClearColor(0.1, 0.1, 0.5, 1.0)
		//gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.Clear(gl.DEPTH_BUFFER_BIT)

		drawScene(grid, window, camera, program)

		// end of draw loop

		grid = wave.NextStep(grid)
	}

	return nil
}

func drawScene(grid wave.Grid, window *win.Window, camera *cam.FpsCamera, program *gfx.Program) {

	point := func(i, j int) mgl32.Vec3 {
		x0 := i
		y0 := grid[2][i][j]
		z0 := j
		return mgl32.Vec3{float32(x0), float32(y0), float32(z0)}
	}

	// p_0 and p_1 clockwise
	normal := func(p_center, p_0, p_1 mgl32.Vec3) mgl32.Vec3 {
		return (p_0.Sub(p_center)).Cross(p_1.Sub(p_center)).Normalize().Mul(-1.0)
	}

	vertices := []float32{}
	for i := 1; i < wave.LEN-1; i++ {
		for j := 0; j < wave.LEN-1; j++ {
			p_center := point(i, j)
			p_right := point(i+1, j)
			p_top := point(i, j+1)
			p_topleft := point(i-1, j+1)

			var n mgl32.Vec3

			// Center/Top/Right
			n = normal(p_center, p_top, p_right)
			vertices = append(vertices,
				p_top.X(), p_top.Y(), p_top.Z(), n.X(), n.Y(), n.Z(),
				p_right.X(), p_right.Y(), p_right.Z(), n.X(), n.Y(), n.Z(),
				p_center.X(), p_center.Y(), p_center.Z(), n.X(), n.Y(), n.Z(),
			)

			// Center/Top/Topleft
			n = normal(p_center, p_topleft, p_top)
			vertices = append(vertices,
				p_top.X(), p_top.Y(), p_top.Z(), n.X(), n.Y(), n.Z(),
				p_topleft.X(), p_topleft.Y(), p_topleft.Z(), n.X(), n.Y(), n.Z(),
				p_center.X(), p_center.Y(), p_center.Z(), n.X(), n.Y(), n.Z(),
			)

		}
	}

	//VAO := createTriangleVAO(vertices, indices)
	VAO, clearVAOFunc := createTriangleVAO(vertices, nil)

	// creates perspective
	fov := float32(60.0)
	projectTransform := mgl32.Perspective(mgl32.DegToRad(fov),
		float32(window.Width())/float32(window.Height()),
		0.01,
		100.0)

	camTransform := camera.GetTransform()
	worldTransform := mgl32.Scale3D(0.2, 0.2, 0.2)
	gl.UniformMatrix4fv(program.GetUniformLocation("project"), 1, false, &projectTransform[0])
	gl.UniformMatrix4fv(program.GetUniformLocation("view"), 1, false, &camTransform[0])
	gl.UniformMatrix4fv(program.GetUniformLocation("world"), 1, false, &worldTransform[0])

	gl.BindVertexArray(VAO)
	//gl.DrawArrays(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, unsafe.Pointer(nil))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)))
	gl.BindVertexArray(0)
	clearVAOFunc()
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action,
	mods glfw.ModifierKey) {

	// When a user presses the escape key, we set the WindowShouldClose property to true,
	// which closes the application
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func interpolation(x1, y1, x2, y2, x3, y3 float64) func(float64) float64 {
	return func(x float64) float64 {
		t1 := y1 * (x - x2) * (x - x3) / ((x1 - x2) * (x1 - x3))
		t2 := y2 * (x - x1) * (x - x3) / ((x2 - x1) * (x2 - x3))
		t3 := y3 * (x - x1) * (x - x2) / ((x3 - x1) * (x3 - x2))
		return t1 + t2 + t3
	}
}
