package gerber_rs274x

import (
	"fmt"
	"math"
	cairo "github.com/ungerik/go-cairo"
)

type MacroAperture struct {
	apertureNumber int
	macroName string
}

func (aperture *MacroAperture) AperturePlaceholder() {

}

func (aperture *MacroAperture) GetHole() Hole {
	return nil
}

func (aperture *MacroAperture) SetHole(hole Hole) {
	
}

func (aperture *MacroAperture) GetMinSize() float64 {
	//TODO: Implement appropriately
	return math.MaxFloat64
}

func (aperture *MacroAperture) DrawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error {
	return nil
}

func (aperture *MacroAperture) DrawApertureSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	return nil
}

func (aperture *MacroAperture) String() string {
	return fmt.Sprintf("{MA, Name: %s}", aperture.macroName)
}