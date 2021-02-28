package gophysics

import "github.com/rpagliuca/go-physics/pkg/gophysics/algebra"

type GravitySource interface {
	GetPotentialEnergy(BodyState) float64
	GetAcceleration(BodyState) Acceleration
	GetX() float64
	GetY() float64
	GetOtherX() float64
	GetOtherY() float64
	GetWidth() float64
	UpdateCenter(x, y int)
	Clone() GravitySource
}

type LinearGravitySource struct {
	Line algebra.Line
}

func (l LinearGravitySource) Clone() GravitySource {
	other := LinearGravitySource{
		l.Line,
	}
	return GravitySource(&other)
}

func (LinearGravitySource) GetPotentialEnergy(BodyState) float64 {
	return 0
}

func (*LinearGravitySource) UpdateCenter(x, y int) {
	// Do nothing
}

func (l LinearGravitySource) GetWidth() float64 {
	return l.Line.Length()
}

func (l LinearGravitySource) GetX() float64 {
	return l.Line.X0
}

func (l LinearGravitySource) GetY() float64 {
	return l.Line.Y0
}

func (l LinearGravitySource) GetOtherX() float64 {
	return l.Line.X1
}

func (l LinearGravitySource) GetOtherY() float64 {
	return l.Line.Y1
}

func (s LinearGravitySource) GetAcceleration(b BodyState) Acceleration {

	normalizedAcceleration, err := algebra.PerpendicularDecomposition(
		s.Line, algebra.Point{b.X, b.Y},
	)

	if err != nil {
		return Acceleration{0, 0}
	}

	// Multiplicar vetor normal pela intensidade da gravidade
	acceleration := Acceleration{GRAVITY * normalizedAcceleration.X1, GRAVITY * normalizedAcceleration.Y1}

	return acceleration
}

type PointGravitySource struct {
	Point algebra.Point
}

func (p *PointGravitySource) Clone() GravitySource {
	other := PointGravitySource{
		p.Point,
	}
	return GravitySource(&other)
}

func (p PointGravitySource) GetPotentialEnergy(bodyState BodyState) float64 {
	return 0
}

func (p PointGravitySource) GetAcceleration(bodyState BodyState) Acceleration {
	return calculateCenterGravity(p.Point, bodyState)
}

func (PointGravitySource) GetWidth() float64 {
	return BOX_SIZE
}

func (p PointGravitySource) GetX() float64 {
	return p.Point.X - BOX_SIZE/2
}

func (p PointGravitySource) GetY() float64 {
	return p.Point.Y
}

func (p PointGravitySource) GetOtherX() float64 {
	return p.Point.X + BOX_SIZE/2
}

func (p PointGravitySource) GetOtherY() float64 {
	return p.Point.Y
}

func (p *PointGravitySource) UpdateCenter(x, y int) {
	p.Point.X = float64(x)
	p.Point.Y = float64(y)
}
