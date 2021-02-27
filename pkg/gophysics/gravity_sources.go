package gophysics

type GravitySource interface {
	getPotentialEnergy(BodyState) float64
	getAcceleration(BodyState) Acceleration
	getX() float64
	getY() float64
	getOtherX() float64
	getOtherY() float64
	getWidth() float64
	updateCenter(x, y int)
}

type LinearGravitySource struct {
	line Line
}

func (LinearGravitySource) getPotentialEnergy(BodyState) float64 {
	return 0
}

func (*LinearGravitySource) updateCenter(x, y int) {
	// Do nothing
}

func (LinearGravitySource) getWidth() float64 {
	return SCREEN_WIDTH
}

func (l LinearGravitySource) getX() float64 {
	return l.line.X0
}

func (l LinearGravitySource) getY() float64 {
	return l.line.Y0
}

func (l LinearGravitySource) getOtherX() float64 {
	return l.line.X1
}

func (l LinearGravitySource) getOtherY() float64 {
	return l.line.Y1
}

func (s LinearGravitySource) getAcceleration(b BodyState) Acceleration {

	normalizedAcceleration, err := perpendicularDecomposition(s.line, Point{b.X, b.Y})

	if err != nil {
		return Acceleration{0, 0}
	}

	// Multiplicar vetor normal pela intensidade da gravidade
	acceleration := Acceleration{GRAVITY * normalizedAcceleration.X1, GRAVITY * normalizedAcceleration.Y1}

	return acceleration
}

type PointGravitySource struct {
	point Point
}

func (p PointGravitySource) getPotentialEnergy(bodyState BodyState) float64 {
	return 0
}

func (p PointGravitySource) getAcceleration(bodyState BodyState) Acceleration {
	return calculateCenterGravity(p.point, bodyState)
}

func (PointGravitySource) getWidth() float64 {
	return BOX_SIZE
}

func (p PointGravitySource) getX() float64 {
	return p.point.X - BOX_SIZE/2
}

func (p PointGravitySource) getY() float64 {
	return p.point.Y
}

func (p PointGravitySource) getOtherX() float64 {
	return p.point.X + BOX_SIZE/2
}

func (p PointGravitySource) getOtherY() float64 {
	return p.point.Y
}

func (p *PointGravitySource) updateCenter(x, y int) {
	p.point.X = float64(x)
	p.point.Y = float64(y)
}
