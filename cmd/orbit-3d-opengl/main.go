package main

/*
http://www.learnopengl.com/#!Lighting/Materials

Shows basic materials with phong lighting
*/
import (
	"log"
	"runtime"

	"github.com/cstegel/opengl-samples-golang/light-maps/cam"
	"github.com/cstegel/opengl-samples-golang/light-maps/gfx"
	"github.com/cstegel/opengl-samples-golang/light-maps/win"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/rpagliuca/go-physics/pkg/algebra"
	"github.com/rpagliuca/go-physics/pkg/dynamics"
)

// vertices to draw 6 faces of a cube
var cubeVertices = []float32{
	// position        // normal vector
	-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
	0.5, -0.5, -0.5, 0.0, 0.0, -1.0,
	0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
	0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
	-0.5, 0.5, -0.5, 0.0, 0.0, -1.0,
	-0.5, -0.5, -0.5, 0.0, 0.0, -1.0,

	-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
	0.5, -0.5, 0.5, 0.0, 0.0, 1.0,
	0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
	0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
	-0.5, 0.5, 0.5, 0.0, 0.0, 1.0,
	-0.5, -0.5, 0.5, 0.0, 0.0, 1.0,

	-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,
	-0.5, 0.5, -0.5, -1.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, -1.0, 0.0, 0.0,
	-0.5, -0.5, 0.5, -1.0, 0.0, 0.0,
	-0.5, 0.5, 0.5, -1.0, 0.0, 0.0,

	0.5, 0.5, 0.5, 1.0, 0.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
	0.5, -0.5, -0.5, 0.0, -1.0, 0.0,
	0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
	0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, -1.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, -1.0, 0.0,

	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
	0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0,
}

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

	window := win.NewWindow(1280, 720, "Lighting maps")

	// Initialize Glow (go function bindings)
	if err := gl.Init(); err != nil {
		panic(err)
	}

	err := programLoop(window)
	if err != nil {
		log.Fatalln(err)
	}
}

/*
 * Creates the Vertex Array Object for a triangle.
 * indices is leftover from earlier samples and not used here.
 */
func createVAO(vertices []float32, indices []uint32) uint32 {

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// size of one whole vertex (sum of attrib sizes)
	var stride int32 = 3*4 + 3*4
	var offset int = 0

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(0)
	offset += 3 * 4

	// normal
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
	gl.EnableVertexAttribArray(1)
	offset += 3 * 4

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO
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

	lightFragShader, err := gfx.NewShaderFromFile("shaders/light.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		return err
	}

	// special shader program so that lights themselves are not affected by lighting
	lightProgram, err := gfx.NewProgram(vertShader, lightFragShader)
	if err != nil {
		return err
	}

	VAO := createVAO(cubeVertices, nil)
	lightVAO := createVAO(cubeVertices, nil)

	// ensure that triangles that are "behind" others do not draw over top of them
	gl.Enable(gl.DEPTH_TEST)

	camera := cam.NewFpsCamera(mgl32.Vec3{20, 60, -40}, mgl32.Vec3{0, 1, 0}, 80, -30, window.InputManager())

	//state := vanillaGravity
	state := multiYinYang

	for !window.ShouldClose() {

		// swaps in last buffer, polls for window events, and generally sets up for a new render frame
		window.StartFrame()

		// update camera position and direction from input evevnts
		camera.Update(window.SinceLastFrame())

		// background color
		//gl.ClearColor(135.0/255.0, 206.0/255.0, 250.0/255.0, 1.0)
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // depth buffer needed for DEPTH_TEST

		// cube rotation matrices
		rotateX := (mgl32.Rotate3DX(mgl32.DegToRad(-45 * float32(glfw.GetTime()))))
		rotateY := (mgl32.Rotate3DY(mgl32.DegToRad(-45 * float32(glfw.GetTime()))))
		rotateZ := (mgl32.Rotate3DZ(mgl32.DegToRad(-45 * float32(glfw.GetTime()))))

		// creates perspective
		fov := float32(60.0)
		projectTransform := mgl32.Perspective(mgl32.DegToRad(fov),
			float32(window.Width())/float32(window.Height()),
			0.1,
			1000.0)

		camTransform := camera.GetTransform()
		lightPos := mgl32.Vec3{0.6, 1, 0.1}
		lightTransform := mgl32.Translate3D(lightPos.X(), lightPos.Y(), lightPos.Z()).Mul4(
			mgl32.Scale3D(0.2, 0.2, 0.2))

		program.Use()
		gl.UniformMatrix4fv(program.GetUniformLocation("view"), 1, false, &camTransform[0])
		gl.UniformMatrix4fv(program.GetUniformLocation("project"), 1, false,
			&projectTransform[0])

		gl.BindVertexArray(VAO)

		// draw each cube after all coordinate system transforms are bound

		// obj is colored, light is white
		gl.Uniform3f(program.GetUniformLocation("material.ambient"), 0.0, 1.0, 0.0)
		gl.Uniform3f(program.GetUniformLocation("material.diffuse"), 0.0, 1.0, 0.0)
		gl.Uniform3f(program.GetUniformLocation("material.specular"), 0.0, 1.0, 0.0)
		gl.Uniform1f(program.GetUniformLocation("material.shininess"), 1.0)

		/*
			lightColor := mgl32.Vec3{
				float32(math.Sin(glfw.GetTime() * 1)),
				float32(math.Sin(glfw.GetTime() * 0.35)),
				float32(math.Sin(glfw.GetTime() * 0.65)),
			}
		*/
		lightColor := mgl32.Vec3{1.0, 1.0, 1.0}

		diffuseColor := mgl32.Vec3{
			0.5 * lightColor[0],
			0.5 * lightColor[1],
			0.5 * lightColor[2],
		}
		ambientColor := mgl32.Vec3{
			0.2 * lightColor[0],
			0.2 * lightColor[1],
			0.2 * lightColor[2],
		}

		gl.Uniform3f(program.GetUniformLocation("light.ambient"),
			ambientColor[0], ambientColor[1], ambientColor[2])
		gl.Uniform3f(program.GetUniformLocation("light.diffuse"),
			diffuseColor[0], diffuseColor[1], diffuseColor[2])
		gl.Uniform3f(program.GetUniformLocation("light.specular"), 1.0, 1.0, 1.0)
		gl.Uniform3f(program.GetUniformLocation("light.position"), lightPos.X(), lightPos.Y(), lightPos.Z())

		state = dynamics.UpdateState(state)

		bodyPositions := [][]float32{}
		for _, body := range state.Bodies {
			bodyPositions = append(bodyPositions, []float32{float32(body.X) / 10.0, float32(body.Y) / 10.0, 0.0})
		}
		for _, gravitySources := range state.GravitySources {
			bodyPositions = append(bodyPositions, []float32{float32(gravitySources.GetX() / 10.0), float32(gravitySources.GetY() / 10.0), 0.0})
		}

		for _, pos := range bodyPositions {

			// turn the cubes into rectangular prisms for more fun
			worldTranslate := mgl32.Translate3D(pos[0], pos[1], pos[2])
			_ = worldTranslate.Mul4(
				rotateX.Mul3(rotateY).Mul3(rotateZ).Mat4(),
			)

			gl.UniformMatrix4fv(program.GetUniformLocation("model"), 1, false,
				&worldTranslate[0])

			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}
		gl.BindVertexArray(0)

		// Draw the light obj after the other boxes using its separate shader program
		// this means that we must re-bind any uniforms
		lightProgram.Use()
		gl.BindVertexArray(lightVAO)
		gl.UniformMatrix4fv(lightProgram.GetUniformLocation("model"), 1, false, &lightTransform[0])
		gl.UniformMatrix4fv(lightProgram.GetUniformLocation("view"), 1, false, &camTransform[0])
		gl.UniformMatrix4fv(lightProgram.GetUniformLocation("project"), 1, false, &projectTransform[0])
		//gl.DrawArrays(gl.TRIANGLES, 0, 36)

		gl.BindVertexArray(0)

		// end of draw loop
	}

	return nil
}

const PIXELS_PER_METER = 10.0
const FRAME_RATE = 60.0
const BOX_SIZE = 10.0
const LINE_WIDTH = 10.0
const SCREEN_WIDTH = 500.0
const SCREEN_HEIGHT = 500.0

var SETTINGS = dynamics.Settings{
	ViewportWidth:       SCREEN_WIDTH,
	ViewportHeight:      SCREEN_HEIGHT,
	ViewportBoxSize:     BOX_SIZE,
	GravityAcceleration: 9.8 * PIXELS_PER_METER,
	DeltaTime:           1.0 / FRAME_RATE,
}

var multiYinYang = dynamics.State{
	SETTINGS,
	[]dynamics.BodyState{
		{X: SCREEN_WIDTH / 2.0, Y: SCREEN_HEIGHT / 20.0 * 1.0, VX: 15.0 * PIXELS_PER_METER, VY: 0.0},
		{X: SCREEN_WIDTH / 2.0, Y: SCREEN_HEIGHT / 20.0 * 3.0, VX: 15.0 * PIXELS_PER_METER, VY: 0.0},
		{X: SCREEN_WIDTH / 2.0, Y: SCREEN_HEIGHT / 20.0 * 5.0, VX: 15.0 * PIXELS_PER_METER, VY: 0.0},
		{X: SCREEN_WIDTH / 2.0, Y: SCREEN_HEIGHT / 20.0 * 7.0, VX: 15.0 * PIXELS_PER_METER, VY: 0.0},
		{X: SCREEN_WIDTH / 2.0, Y: SCREEN_HEIGHT / 20.0 * 13.0, VX: -15.0 * PIXELS_PER_METER, VY: 0.0},
		{X: SCREEN_WIDTH / 2.0, Y: SCREEN_HEIGHT / 20.0 * 15.0, VX: -15.0 * PIXELS_PER_METER, VY: 0.0},
		{X: SCREEN_WIDTH / 2.0, Y: SCREEN_HEIGHT / 20.0 * 17.0, VX: -15.0 * PIXELS_PER_METER, VY: 0.0},
		{X: SCREEN_WIDTH / 2.0, Y: SCREEN_HEIGHT / 20.0 * 19.0, VX: -15.0 * PIXELS_PER_METER, VY: 0.0},
	},
	[]dynamics.GravitySource{
		// Fonte gravitacional pontual, como se fosse um movimento astronômico
		&dynamics.PointGravitySource{SETTINGS, algebra.Point{SCREEN_WIDTH / 2.0, SCREEN_HEIGHT / 2.0}},
	},
}

var vanillaGravity = dynamics.State{
	SETTINGS,
	[]dynamics.BodyState{
		// 3 corpos
		{X: 0.0, Y: 0.0, VX: 15.0 * PIXELS_PER_METER, VY: 3.0 * PIXELS_PER_METER},
		{X: SCREEN_WIDTH / 2.0, Y: 80.0, VX: -40.0 * PIXELS_PER_METER, VY: -1.0 * PIXELS_PER_METER},
		{X: SCREEN_WIDTH / 2.0, Y: 200.0, VX: -30.0 * PIXELS_PER_METER, VY: -7.0 * PIXELS_PER_METER},
	},
	[]dynamics.GravitySource{
		// 1 fonte gravitacional no chão (gravidade padrão, como estamos acostumados)
		&dynamics.LinearGravitySource{SETTINGS, algebra.Line{0.0, SCREEN_HEIGHT, SCREEN_WIDTH, SCREEN_HEIGHT}}, // Bottom
	},
}
