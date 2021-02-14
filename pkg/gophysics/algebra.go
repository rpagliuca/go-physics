package gophysics

import "math"

func perpendicularLine(l Line) Line {
	//Rotation by 90º
	// | 0 -1 | (X) = (Xf)
	// | 1  0 | (Y) = (Yf)
	return Line{-l.Y0, l.X0, -l.Y1, l.X1}
}

func normalizeLine(l Line) Line {
	magnitude := math.Pow(math.Pow(l.X1-l.X0, 2)+math.Pow(l.Y1-l.Y0, 2), 0.5)
	return Line{0, 0, (l.X1 - l.X0) / magnitude, (l.Y1 - l.Y0) / magnitude}
}

func perpendicularDecomposition(line Line, point Point) Line {
	x := 0.0
	y := 0.0

	if line.coefficientA() == 0 { // Horizontal gravity source
		x = point.X
		y = line.Y0
	} else if math.IsInf(line.coefficientA(), 0) { // Vertical gravity source
		x = line.X0
		y = point.Y
	} else { // Generic case
		// Calcular uma reta qualquer perpendicular à gravitySource
		p := normalizeLine(perpendicularLine(line))

		// Calcular uma nova reta, também perpendicular, mas passando pelo body
		r1 := Line{point.X + p.X1, point.Y + p.Y1, point.X, point.Y}

		// Encontrar pontos X e Y da interseção
		deltaB := line.coefficientB() - r1.coefficientB()
		deltaA := r1.coefficientA() - line.coefficientA()

		x = deltaB / deltaA
		y = r1.calculateY(x)
	}

	return normalizeLine(Line{point.X, point.Y, x, y})
}

type Line struct {
	X0, Y0, X1, Y1 float64
}

func (l Line) coefficientA() float64 {
	l2 := normalizeLine(l)
	a := l2.Y1 / l2.X1
	return a
}

func (l Line) coefficientB() float64 {
	return l.Y0 - l.coefficientA()*l.X0
}

func (l Line) calculateY(x float64) float64 {
	return x*l.coefficientA() + l.coefficientB()
}

type Point struct {
	X, Y float64
}
