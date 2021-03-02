package dynamics

const BOUNCING_CONSERVATION = 0.3

func getAcceleration(bodyState BodyState, gravitySources []GravitySource) Acceleration {
	acceleration := Acceleration{0, 0}
	for i := range gravitySources {
		newAcceleration := gravitySources[i].GetAcceleration(bodyState)
		acceleration.AX += newAcceleration.AX
		acceleration.AY += newAcceleration.AY
	}
	return acceleration
}

func UpdateState(state State) State {
	for i := range state.Bodies {
		bodyState := state.Bodies[i]
		acceleration := getAcceleration(bodyState, state.GravitySources)
		nextState := getNextBodyStateRungeKutta(bodyState, acceleration, state.Settings)
		state.Bodies[i] = nextState
	}
	return state
}
