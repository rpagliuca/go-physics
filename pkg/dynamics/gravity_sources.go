package dynamics

import "github.com/rpagliuca/go-physics/pkg/algebra"

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
	Settings Settings
	Line     algebra.Line
}

func (l LinearGravitySource) Clone() GravitySource {
	other := LinearGravitySource{
		l.Settings.Clone(),
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
	normalized, err := algebra.PerpendicularDecomposition(
		s.Line, algebra.Point{b.X, b.Y},
	)
	if err != nil {
		return Acceleration{0, 0}
	}
	acc := Acceleration{
		s.Settings.GravityAcceleration * normalized.X1,
		s.Settings.GravityAcceleration * normalized.Y1,
	}
	return acc
}

type PointGravitySource struct {
	Settings Settings
	Point    algebra.Point
}

func (p *PointGravitySource) Clone() GravitySource {
	other := PointGravitySource{
		p.Settings.Clone(),
		p.Point,
	}
	return GravitySource(&other)
}

func (p PointGravitySource) GetPotentialEnergy(bodyState BodyState) float64 {
	return 0
}

func (p PointGravitySource) GetAcceleration(bodyState BodyState) Acceleration {
	vector := algebra.Line{bodyState.X, bodyState.Y, p.Point.X, p.Point.Y}
	normalized := algebra.NormalizeLine(vector)
	acc := Acceleration{p.Settings.GravityAcceleration * normalized.X1, p.Settings.GravityAcceleration * normalized.Y1}
	return acc
}

func (p PointGravitySource) GetWidth() float64 {
	return p.Settings.ViewportBoxSize
}

func (p PointGravitySource) GetX() float64 {
	return p.Point.X - p.Settings.ViewportBoxSize/2
}

func (p PointGravitySource) GetY() float64 {
	return p.Point.Y
}

func (p PointGravitySource) GetOtherX() float64 {
	return p.Point.X + p.Settings.ViewportBoxSize/2
}

func (p PointGravitySource) GetOtherY() float64 {
	return p.Point.Y
}

func (p *PointGravitySource) UpdateCenter(x, y int) {
	p.Point.X = float64(x)
	p.Point.Y = float64(y)
}
