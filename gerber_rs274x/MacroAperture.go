package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

type MacroAperture struct {
	macroName string
}

func (macro* MacroAperture) AperturePlaceholder() {

}

func (macro* MacroAperture) GetHole() Hole {
	return nil
}

func (macro* MacroAperture) SetHole(hole Hole) {
	
}

func (macro* MacroAperture) DrawApertureSVG(svg *svg.SVG, gfxState *GraphicsState, x float64, y float64) error {
	return nil
}

func (macro* MacroAperture) String() string {
	return fmt.Sprintf("{MA, Name: %s}", macro.macroName)
}