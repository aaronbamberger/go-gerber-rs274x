package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

type ObroundAperture struct {
	xSize float64
	ySize float64
	Hole
}

func (obround* ObroundAperture) AperturePlaceholder() {

}

func (obround* ObroundAperture) GetHole() Hole {
	return obround.Hole
}

func (obround* ObroundAperture) SetHole(hole Hole) {
	obround.Hole = hole
}

func (obround* ObroundAperture) DrawApertureSVG(svg *svg.SVG, gfxState *GraphicsState, x float64, y float64) error {
	return nil
}

func (obround* ObroundAperture) String() string {
	return fmt.Sprintf("{OA, X: %f, Y: %f, Hole: %v}", obround.xSize, obround.ySize, obround.Hole)
}