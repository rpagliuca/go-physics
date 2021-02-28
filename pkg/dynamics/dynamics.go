package dynamics

import "fmt"

func getAcceleration(bodyState BodyState, gravitySources []GravitySource) Acceleration {
	acceleration := Acceleration{0, 0}
	for i := range gravitySources {
		newAcceleration := gravitySources[i].GetAcceleration(bodyState)
		acceleration.AX += newAcceleration.AX
		acceleration.AY += newAcceleration.AY
	}
	return acceleration
}

func getNextBodyState(state BodyState, acceleration Acceleration, settings Settings) BodyState {
	// TODO implement Runge-Kutta
	nextBodyState := BodyState{}
	nextBodyState.VX = state.VX + acceleration.AX*settings.DeltaTime
	nextBodyState.VY = state.VY + acceleration.AY*settings.DeltaTime
	nextBodyState.X = state.X + (state.VX+nextBodyState.VX)/2
	nextBodyState.Y = state.Y + (state.VY+nextBodyState.VY)/2
	if nextBodyState.Y > settings.ViewportHeight-settings.ViewportBoxSize {
		nextBodyState.Y = settings.ViewportHeight - settings.ViewportBoxSize
		nextBodyState.VY = -0.95 * state.VY
	}
	if nextBodyState.Y < 0 {
		nextBodyState.Y = 0
		nextBodyState.VY = -0.95 * state.VY
	}
	if nextBodyState.X > settings.ViewportWidth-settings.ViewportBoxSize {
		nextBodyState.X = settings.ViewportWidth - settings.ViewportBoxSize
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
		fmt.Println("bodyState", bodyState)
		acceleration := getAcceleration(bodyState, state.GravitySources)
		fmt.Println("acceleration", acceleration)
		nextState := getNextBodyState(bodyState, acceleration, state.Settings)
		fmt.Println("nextState", nextState)
		state.Bodies[i] = nextState
	}
	return state
}
