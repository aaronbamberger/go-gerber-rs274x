package gerber_rs274x

import "fmt"

type CircleAperture struct {
	diameter float64
	Hole
}

func (circle* CircleAperture) AperturePlaceholder() {

}

func (circle* CircleAperture) GetHole() Hole {
	return circle.Hole
}

func (circle* CircleAperture) SetHole(hole Hole) {
	circle.Hole = hole
}

func (circle* CircleAperture) String() string {
	return fmt.Sprintf("{CA, Diameter: %f, Hole: %v}", circle.diameter, circle.Hole)
}