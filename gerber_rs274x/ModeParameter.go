package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

type ModeParameter struct {
	paramCode ParameterCode
	units Units
}

func (mode *ModeParameter) DataBlockPlaceholder() {

}

func (mode *ModeParameter) ProcessDataBlockSVG(svg *svg.SVG, gfxState *GraphicsState) error {
	//TODO: For now this doesn't alter the graphics state or draw anything
	return nil
}

func (moParam *ModeParameter) String() string {
	var units string
	
	switch moParam.units {
		case UNITS_IN:
			units = "Inches"
			
		case UNITS_MM:
			units = "Millimeters"
			
		default:
			units = "Unknown"
	}
	
	return fmt.Sprintf("{MO, Units: %s}", units)
}