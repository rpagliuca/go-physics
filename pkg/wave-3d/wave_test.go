package wave3d

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWave(t *testing.T) {

	g0 := Grid{}

	g0[0][0][0] = 100
	g0[0][1][1] = 90
	g0[0][2][2] = 80
	g0[0][3][3] = 70

	g0[1] = g0[0]
	g0[2] = g0[0]

	g1 := NextStep(g0)

	assert.NotEqual(t, g0, g1)
}
