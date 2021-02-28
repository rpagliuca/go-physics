package dynamics

func saveGlobals() []float64 {
	return []float64{
		GRAVITY,
		FRAME_RATE,
		SCREEN_WIDTH,
		SCREEN_HEIGHT,
		BOX_SIZE,
	}
}

func restoreGlobals(globals []float64) {
	GRAVITY = globals[0]
	FRAME_RATE = globals[1]
	SCREEN_WIDTH = globals[2]
	SCREEN_HEIGHT = globals[3]
	BOX_SIZE = globals[4]
}
