package gophysics

import (
	"errors"
	"math"
)

func PerpendicularLine(l Line) Line {
	//Rotation by 90º
	// | 0 -1 | (X) = (Xf)
	// | 1  0 | (Y) = (Yf)
	return Line{-l.Y0, l.X0, -l.Y1, l.X1}
}

func NormalizeLine(l Line) Line {
	magnitude := l.Length()
	return Line{0, 0, (l.X1 - l.X0) / magnitude, (l.Y1 - l.Y0) / magnitude}
}

func PerpendicularDecomposition(line Line, point Point) (Line, error) {
	x := 0.0
	y := 0.0

	if line.CoefficientA() == 0 { // Horizontal gravity source
		x = point.X
		y = line.Y0
	} else if math.IsInf(line.CoefficientA(), 0) { // Vertical gravity source
		x = line.X0
		y = point.Y
	} else { // Generic case
		// Calcular uma reta qualquer perpendicular à gravitySource
		p := NormalizeLine(PerpendicularLine(line))

		// Calcular uma nova reta, também perpendicular, mas passando pelo body
		r1 := Line{point.X + p.X1, point.Y + p.Y1, point.X, point.Y}

		// Encontrar pontos X e Y da interseção
		deltaB := line.CoefficientB() - r1.CoefficientB()
		deltaA := r1.CoefficientA() - line.CoefficientA()

		x = deltaB / deltaA
		y = r1.CalculateY(x)
	}

	normalized := NormalizeLine(Line{point.X, point.Y, x, y})

	// Edge case
	if math.IsNaN(normalized.X1) || math.IsNaN(normalized.Y1) {
		return Line{0, 0, 0, 0}, errors.New("Failed calculating perpendicular line from point to line. Perhaps the point is on the line.")
	}

	return normalized, nil
}

type Line struct {
	X0, Y0, X1, Y1 float64
}

func (l Line) CoefficientA() float64 {
	l2 := NormalizeLine(l)
	a := l2.Y1 / l2.X1
	return a
}

func (l Line) CoefficientB() float64 {
	return l.Y0 - l.CoefficientA()*l.X0
}

func (l Line) CalculateY(x float64) float64 {
	return x*l.CoefficientA() + l.CoefficientB()
}

func (l Line) Length() float64 {
	return math.Pow(math.Pow(l.X1-l.X0, 2)+math.Pow(l.Y1-l.Y0, 2), 0.5)
}

type Point struct {
	X, Y float64
}
