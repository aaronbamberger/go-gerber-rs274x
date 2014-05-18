package gerber_rs274x

import (
	"fmt"
	cairo "github.com/ungerik/go-cairo"
)

type RectangularHole struct {
	holeXSize float64
	holeYSize float64
}

func (rectangle* RectangularHole) HolePlaceholder() {

}

func (rectangle* RectangularHole) String() string {
	return fmt.Sprintf("{RH, X: %f, Y: %f}", rectangle.holeXSize, rectangle.holeYSize)
}

func (hole *RectangularHole) DrawHoleSurface(surface *cairo.Surface) error {
	
	xRadius := hole.holeXSize / 2.0
	yRadius := hole.holeYSize / 2.0
	
	surface.Save()
	
	// We temporarily set the compositing operator to clear, to clear the hole to transparent
	surface.SetOperator(cairo.OPERATOR_CLEAR)
	surface.MoveTo(-xRadius, -yRadius)
	surface.LineTo(xRadius, -yRadius)
	surface.LineTo(xRadius, yRadius)
	surface.LineTo(-xRadius, yRadius)
	surface.LineTo(-xRadius, -yRadius)
	surface.Fill()
	
	surface.Restore()
	
	return nil
}