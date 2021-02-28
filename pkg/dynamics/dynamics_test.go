package dynamics

import (
	"testing"

	"github.com/rpagliuca/go-physics/pkg/algebra"
	"github.com/stretchr/testify/assert"
)

var SETTINGS = Settings{
	GravityAcceleration: 1,
	DeltaTime:           1,
	ViewportHeight:      1000,
	ViewportWidth:       1000,
	ViewportBoxSize:     10,
}

func TestLinearGravitySource(t *testing.T) {

	s0 := State{
		SETTINGS,
		[]BodyState{
			{
				X:  5,
				Y:  10,
				VX: 1,
				VY: 0,
			},
		},
		[]GravitySource{
			&LinearGravitySource{SETTINGS, algebra.Line{0, 0, 10, 0}},
		},
	}

	s1 := UpdateState(s0.Clone())

	assert.Equal(t, 6.0, s1.Bodies[0].X)
	assert.Equal(t, -1.0, s1.Bodies[0].VY)
}

func TestPointGravitySource(t *testing.T) {

	s0 := State{
		SETTINGS,
		[]BodyState{
			{
				X:  5,
				Y:  10,
				VX: 1,
				VY: 0,
			},
		},
		[]GravitySource{
			&PointGravitySource{SETTINGS, algebra.Point{5, 5}},
		},
	}

	s1 := UpdateState(s0.Clone())

	assert.Equal(t, 6.0, s1.Bodies[0].X)
	assert.Equal(t, -1.0, s1.Bodies[0].VY)

	s2 := UpdateState(s1.Clone())

	assert.NotEqual(t, 7.0, s2.Bodies[0].X)
	assert.Less(t, s2.Bodies[0].VY, -1.0)
}
