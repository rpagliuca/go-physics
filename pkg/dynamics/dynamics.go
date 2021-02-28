package dynamics

// TODO: Remove constant frame rate and gravity
var FRAME_RATE = 60.0
var GRAVITY = 9.8

// TODO: Remove dependency on camera/viewport
var BOX_SIZE = 10.0
var SCREEN_WIDTH = 320.0
var SCREEN_HEIGHT = 240.0

func fixAccelerationRate(a Acceleration, frameRate float64) Acceleration {
	return Acceleration{a.AX / frameRate, a.AY / frameRate}
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
		nextState := getNextBodyState(bodyState, fixAccelerationRate(acceleration, FRAME_RATE))
		state.Bodies[i] = nextState
	}
	return state
}
