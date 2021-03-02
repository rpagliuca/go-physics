package dynamics

var firstIteration = true

func getNextBodyStateRungeKutta(state BodyState, acceleration Acceleration, settings Settings) BodyState {

	nextBodyState := BodyState{}

	dvx := func(t, vx float64) float64 {
		return acceleration.AX
	}
	nextBodyState.VX = rungeKutta(state.VX, settings.DeltaTime, dvx)

	dvy := func(t, vy float64) float64 {
		return acceleration.AY
	}
	nextBodyState.VY = rungeKutta(state.VY, settings.DeltaTime, dvy)

	dx := func(t, x float64) float64 {
		return state.VX + acceleration.AX*t
	}
	nextBodyState.X = rungeKutta(state.X, settings.DeltaTime, dx)

	dy := func(t, y float64) float64 {
		return state.VY + acceleration.AY*t
	}
	nextBodyState.Y = rungeKutta(state.Y, settings.DeltaTime, dy)

	if nextBodyState.Y > settings.ViewportHeight-settings.ViewportBoxSize {
		nextBodyState.Y = settings.ViewportHeight - settings.ViewportBoxSize
		nextBodyState.VY = -BOUNCING_CONSERVATION * state.VY
		nextBodyState.VX = 0.95 * nextBodyState.VX
	}
	if nextBodyState.Y < 0 {
		nextBodyState.Y = 0
		nextBodyState.VY = -BOUNCING_CONSERVATION * state.VY
		nextBodyState.VX = 0.95 * nextBodyState.VX
	}
	if nextBodyState.X > settings.ViewportWidth-settings.ViewportBoxSize {
		nextBodyState.X = settings.ViewportWidth - settings.ViewportBoxSize
		nextBodyState.VX = -BOUNCING_CONSERVATION * state.VX
		nextBodyState.VY = 0.95 * nextBodyState.VY
	}
	if nextBodyState.X < 0 {
		nextBodyState.X = 0
		nextBodyState.VX = -BOUNCING_CONSERVATION * state.VX
		nextBodyState.VY = 0.95 * nextBodyState.VY
	}

	nextBodyState.VX = 0.9995 * nextBodyState.VX
	nextBodyState.VY = 0.9995 * nextBodyState.VY

	return nextBodyState
}

func rungeKutta(yn float64, deltaTime float64, f func(t, y float64) float64) float64 {
	k1 := f(0, yn)
	k2 := f(deltaTime/2, yn+deltaTime/2*k1)
	k3 := f(deltaTime/2, yn+deltaTime/2*k2)
	k4 := f(deltaTime, yn+deltaTime*k3)
	return yn + deltaTime/6*(k1+2*k2+2*k3+k4)
}
