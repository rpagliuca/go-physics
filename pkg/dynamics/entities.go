package dynamics

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

func (b BodyState) Clone() BodyState {
	return b
}

type State struct {
	Settings       Settings
	Bodies         []BodyState
	GravitySources []GravitySource
}

func (s State) Clone() State {
	bodies := []BodyState{}
	for i := range s.Bodies {
		bodies = append(bodies, s.Bodies[i].Clone())
	}
	gravitySources := []GravitySource{}
	for i := range s.GravitySources {
		gravitySources = append(gravitySources, s.GravitySources[i].Clone())
	}
	return State{
		s.Settings.Clone(),
		bodies,
		gravitySources,
	}
}

type Settings struct {
	ViewportWidth       float64
	ViewportHeight      float64
	ViewportBoxSize     float64
	GravityAcceleration float64
	DeltaTime           float64
}

func (s Settings) Clone() Settings {
	return s
}
