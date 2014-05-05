package gerber_rs274x

import (
	"fmt"
	cairo "github.com/ungerik/go-cairo"
)

type ModeParameter struct {
	paramCode ParameterCode
	units Units
}

func (mode *ModeParameter) DataBlockPlaceholder() {

}

func (mode *ModeParameter) ProcessDataBlockBoundsCheck(imageBounds *ImageBounds, gfxState *GraphicsState) error {
	//TODO: For now this doesn't alter the graphics state or draw anything
	return nil
}

func (mode *ModeParameter) ProcessDataBlockSurface(surface *cairo.Surface, gfxState *GraphicsState) error {
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