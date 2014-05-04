package gerber_rs274x

import "fmt"

type RectangularHole struct {
	holeXSize float64
	holeYSize float64
}

func (rectangle* RectangularHole) HolePlaceholder() {

}

func (rectangle* RectangularHole) String() string {
	return fmt.Sprintf("{RH, X: %f, Y: %f}", rectangle.holeXSize, rectangle.holeYSize)
}