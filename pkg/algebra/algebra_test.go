package algebra

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHorizontalLinePerpendicularDecomposition(t *testing.T) {
	// Reference lines
	right := Line{0, 0, 1, 0}
	down := Line{0, 0, 0, -1}
	up := Line{0, 0, 0, 1}

	// Test cases
	cases := []struct {
		point    Point
		expected Line
	}{
		{
			Point{0, 1},
			down,
		},
		{
			Point{1, 1},
			down,
		},
		{
			Point{0, -1},
			up,
		},
		{
			Point{-1, -1},
			up,
		},
	}

	// Test horizontal line
	for _, c := range cases {
		got, err := PerpendicularDecomposition(right, c.point)
		assert.Nil(t, err)
		assert.Equal(t, got, c.expected)
	}
}

func TestVerticalLinePerpendicularDecomposition(t *testing.T) {
	// Reference lines
	right := Line{0, 0, 1, 0}
	left := Line{0, 0, -1, 0}
	up := Line{0, 0, 0, 1}

	// Test cases
	cases := []struct {
		point    Point
		expected Line
	}{
		{
			Point{1, 0},
			left,
		},
		{
			Point{1, 1},
			left,
		},
		{
			Point{-1, 0},
			right,
		},
		{
			Point{-1, -1},
			right,
		},
	}

	// Test vertical line
	for _, c := range cases {
		got, err := PerpendicularDecomposition(up, c.point)
		assert.Nil(t, err)
		assert.Equal(t, got, c.expected)
	}
}

func TestDiagonalLinePerpendicularDecomposition(t *testing.T) {
	// Reference lines
	diagonal := Line{0, 0, 1, 1}
	diagonal90 := Line{0, 0, 1, -1}
	diagonal270 := Line{0, 0, -1, 1}

	// Test cases
	cases := []struct {
		point    Point
		expected Line
	}{
		{
			Point{1, 0},
			diagonal270,
		},
		{
			Point{0, 1},
			diagonal90,
		},
	}

	// Test diagonal line
	for _, c := range cases {
		got, err := PerpendicularDecomposition(diagonal, c.point)
		assert.Nil(t, err)
		assert.Equal(t, got, NormalizeLine(c.expected))
	}
}

func TestOverlappingLineAndPoint(t *testing.T) {
	// Reference lines
	right := Line{0, 0, 1, 0}

	// Test cases
	cases := []struct {
		point Point
	}{
		{
			Point{1, 0},
		},
		{
			Point{2, 0},
		},
	}

	// Test diagonal line
	for _, c := range cases {
		got, err := PerpendicularDecomposition(right, c.point)
		assert.Equal(t, got, Line{0, 0, 0, 0})
		assert.NotNil(t, err)
	}
}
