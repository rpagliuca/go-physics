package wave

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWave(t *testing.T) {

	p0 := Position{}
	p0[0] = 100
	p0[1] = 90
	p0[2] = 80
	p0[3] = 70

	s0 := State{Position: p0}

	s1 := NextStep(s0)

	assert.NotEqual(t, s0, s1)
}
