package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

type PolygonAperture struct {
	outerDiameter float64
	numVertices int
	rotationDegrees float64
	Hole
}

func (polygon* PolygonAperture) AperturePlaceholder() {

}

func (polygon* PolygonAperture) GetHole() Hole {
	return polygon.Hole
}

func (polygon* PolygonAperture) SetHole(hole Hole) {
	polygon.Hole = hole
}

func (polygon* PolygonAperture) DrawApertureSVG(svg *svg.SVG, gfxState *GraphicsState, x float64, y float64) error {
	return nil
}

func (polygon* PolygonAperture) String() string {
	return fmt.Sprintf("{PA, Diameter: %f, Vertices: %d, Rotation: %f, Hole: %v", polygon.outerDiameter, polygon.numVertices, polygon.rotationDegrees, polygon.Hole)
}