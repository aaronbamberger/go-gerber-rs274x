package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

type ApertureDefinitionParameter struct {
	paramCode ParameterCode
	apertureNumber int
	apertureType ApertureType
	aperture Aperture
}

func (apertureDefinition *ApertureDefinitionParameter) DataBlockPlaceholder() {

}

func (apertureDefinition *ApertureDefinitionParameter) ProcessDataBlockSVG(svg *svg.SVG, gfxState *GraphicsState) error {
	// Remember this aperture in the graphics state for later use
	gfxState.apertures[apertureDefinition.apertureNumber] = apertureDefinition.aperture
	
	return nil
}

func (adParam *ApertureDefinitionParameter) String() string {
	var apertureType string
	
	switch adParam.apertureType {
		case CIRCLE_APERTURE:
			apertureType = "Circle"
		
		case RECTANGLE_APERTURE:
			apertureType = "Rectangle"
		
		case OBROUND_APERTURE:
			apertureType = "Obround"
		
		case POLYGON_APERTURE:
			apertureType = "Polygon"
		
		case MACRO_APERTURE:
			apertureType = "Macro"
		
		default:
			apertureType = "Unknown"
	}

	return fmt.Sprintf("{AD, D-Code: %d, Type: %s, Aperture: %s}", adParam.apertureNumber, apertureType, adParam.aperture)
}