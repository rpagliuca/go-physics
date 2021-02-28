package gophysics

import (
	"testing"

	"github.com/rpagliuca/go-physics/pkg/gophysics/algebra"
	"github.com/stretchr/testify/assert"
)

func TestLinearGravitySource(t *testing.T) {

	globals := saveGlobals()

	GRAVITY = 1
	FRAME_RATE = 1

	s0 := State{
		10,
		10,
		[]BodyState{
			{
				X:  5,
				Y:  10,
				VX: 1,
				VY: 0,
			},
		},
		[]GravitySource{
			&LinearGravitySource{algebra.Line{0, 0, 10, 0}},
		},
	}

	s1 := UpdateState(s0.Clone())

	assert.Equal(t, 6.0, s1.Bodies[0].X)
	assert.Equal(t, -0.99, s1.Bodies[0].VY)

	restoreGlobals(globals)
}

func TestPointGravitySource(t *testing.T) {

	globals := saveGlobals()

	GRAVITY = 1
	FRAME_RATE = 1

	s0 := State{
		10,
		10,
		[]BodyState{
			{
				X:  5,
				Y:  10,
				VX: 1,
				VY: 0,
			},
		},
		[]GravitySource{
			&PointGravitySource{algebra.Point{5, 5}},
		},
	}

	s1 := UpdateState(s0.Clone())

	assert.Equal(t, 6.0, s1.Bodies[0].X)
	assert.Equal(t, -0.99, s1.Bodies[0].VY)

	s2 := UpdateState(s1.Clone())

	assert.NotEqual(t, 7.0, s2.Bodies[0].X)
	assert.Less(t, s2.Bodies[0].VY, -0.99)

	restoreGlobals(globals)
}
