package gerber_rs274x

import (
	"fmt"
	"math"
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

func (hole *CircularHole) DrawHoleSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	
	radius := (hole.holeDiameter / 2.0) * gfxState.scaleFactor
	
	surface.Save()
	
	// We temporarily set the compositing operator to clear, to clear the hole to transparent
	surface.SetOperator(cairo.OPERATOR_CLEAR)
	surface.Arc(x, y, radius, 0, 2.0 * math.Pi)
	surface.Fill()
	
	surface.Restore()
	
	return nil
}