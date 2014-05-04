package gerber_rs274x

import "fmt"

type RectangleAperture struct {
	xSize float64
	ySize float64
	Hole
}

func (rectangle* RectangleAperture) AperturePlaceholder() {

}

func (rectangle* RectangleAperture) GetHole() Hole {
	return rectangle.Hole
}

func (rectangle* RectangleAperture) SetHole(hole Hole) {
	rectangle.Hole = hole
}

func (rectangle* RectangleAperture) String() string {
	return fmt.Sprintf("{RA, X: %f, Y: %f, Hole: %v}", rectangle.xSize, rectangle.ySize, rectangle.Hole)
}