package gophysics

import (
	"math"
)

const GRAVITY = 9.8

// TODO: Remove dependency on camera/viewport
const BOX_SIZE = 10
const SCREEN_WIDTH = 320
const SCREEN_HEIGHT = 240

// TODO: Remove constant frame rate
func fixAccelerationRate(a Acceleration) Acceleration {
	return Acceleration{a.AX / 60, a.AY / 60}
}

func getAcceleration(bodyState BodyState, gravitySources []GravitySource) Acceleration {
	acceleration := Acceleration{0, 0}
	for i := range gravitySources {
		newAcceleration := gravitySources[i].GetAcceleration(bodyState)
		acceleration.AX += newAcceleration.AX
		acceleration.AY += newAcceleration.AY
	}
	return acceleration
}

func calculateCenterGravity(point Point, state BodyState) Acceleration {
	// TODO merge this code with generic line gravity source
	center := []float64{point.X, point.Y}
	pos := []float64{state.X, state.Y}
	direction := []float64{center[0] - pos[0], center[1] - pos[1]}
	magnitude := math.Pow(math.Pow(direction[0], 2)+math.Pow(direction[1], 2), 0.5)
	direction_normalized := []float64{direction[0] / magnitude, direction[1] / magnitude}
	acc := Acceleration{
		GRAVITY * direction_normalized[0],
		GRAVITY * direction_normalized[1],
	}
	return acc
}

func getNextBodyState(state BodyState, acceleration Acceleration) BodyState {
	// TODO implement Runge-Kutta
	nextBodyState := BodyState{}
	nextBodyState.VX = state.VX + acceleration.AX
	nextBodyState.VY = state.VY + acceleration.AY
	nextBodyState.X = state.X + (state.VX+nextBodyState.VX)/2
	nextBodyState.Y = state.Y + (state.VY+nextBodyState.VY)/2
	if nextBodyState.Y > SCREEN_HEIGHT-BOX_SIZE {
		nextBodyState.Y = SCREEN_HEIGHT - BOX_SIZE
		nextBodyState.VY = -0.95 * state.VY
	}
	if nextBodyState.Y < 0 {
		nextBodyState.Y = 0
		nextBodyState.VY = -0.95 * state.VY
	}
	if nextBodyState.X > SCREEN_WIDTH-BOX_SIZE {
		nextBodyState.X = SCREEN_WIDTH - BOX_SIZE
		nextBodyState.VX = -0.95 * state.VX
	}
	if nextBodyState.X < 0 {
		nextBodyState.X = 0
		nextBodyState.VX = -0.95 * state.VX
	}
	nextBodyState.VX *= 0.99
	nextBodyState.VY *= 0.99
	return nextBodyState
}

func UpdateState(state State) State {
	for i := range state.Bodies {
		bodyState := state.Bodies[i]
		acceleration := getAcceleration(bodyState, state.GravitySources)
		nextState := getNextBodyState(bodyState, fixAccelerationRate(acceleration))
		state.Bodies[i] = nextState
	}
	return state
}

type Acceleration struct {
	AX float64
	AY float64
}

type BodyState struct {
	X  float64
	Y  float64
	VX float64
	VY float64
}

type State struct {
	ViewportWidth  int
	ViewportHeight int
	Bodies         []BodyState
	GravitySources []GravitySource
}
