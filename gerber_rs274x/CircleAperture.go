package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

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

func (circle* CircleAperture) DrawApertureSVG(svg *svg.SVG, gfxState *GraphicsState, x float64, y float64) error {
	//TODO: Scaling is temporary, need to figure out how to do this the right way
	scaleFactor := 1000.0
	scaledRadius := (circle.diameter / 2.0) * scaleFactor
	scaledX := x * scaleFactor
	scaledY := y * scaleFactor
	var fill string
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		fill = "fill:black"
	} else {
		fill = "fill:white"
	}
	
	svg.Circle(int(scaledX), int(scaledY), int(scaledRadius), fill)
	
	return nil
}

func (circle* CircleAperture) String() string {
	return fmt.Sprintf("{CA, Diameter: %f, Hole: %v}", circle.diameter, circle.Hole)
}