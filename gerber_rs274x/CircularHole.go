package gerber_rs274x

import (
	"fmt"
	cairo "github.com/ungerik/go-cairo"
)

type CircularHole struct {
	holeDiameter float64
}

func (hole *CircularHole) HolePlaceholder() {

}

func (hole *CircularHole) String() string {
	return fmt.Sprintf("{CH, Diameter: %f}", hole.holeDiameter)
}

func (hole *CircularHole) DrawHoleSurface(surface *cairo.Surface) error {
	
	radius := (hole.holeDiameter / 2.0)
	
	surface.Save()
	
	// We temporarily set the compositing operator to clear, to clear the hole to transparent
	surface.SetOperator(cairo.OPERATOR_CLEAR)
	surface.Arc(0.0, 0.0, radius, 0, TWO_PI)
	surface.Fill()
	
	surface.Restore()
	
	return nil
}