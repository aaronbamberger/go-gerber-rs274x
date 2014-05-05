package gerber_rs274x

import (
	"fmt"
	cairo "github.com/ungerik/go-cairo"
)

type GraphicsStateChange struct {
	fnCode FunctionCode
}

func (graphicsStateChange *GraphicsStateChange) DataBlockPlaceholder() {

}

func (graphicsStateChange *GraphicsStateChange) ProcessDataBlockBoundsCheck(imageBounds *ImageBounds, gfxState *GraphicsState) error {
	switch graphicsStateChange.fnCode {
		case SINGLE_QUADRANT_MODE, MULTI_QUADRANT_MODE:
			gfxState.currentQuadrantMode = graphicsStateChange.fnCode
		
		case REGION_MODE_ON:
			gfxState.regionModeOn = true
			
		case REGION_MODE_OFF:
			gfxState.regionModeOn = false
			
		case END_OF_FILE:
			gfxState.fileComplete = true
			
		// For now, we're not going to do anything with any of the other ones
	}
	
	return nil
}

func (graphicsStateChange *GraphicsStateChange) ProcessDataBlockSurface(surface *cairo.Surface, gfxState *GraphicsState) error {
	switch graphicsStateChange.fnCode {
		case SINGLE_QUADRANT_MODE, MULTI_QUADRANT_MODE:
			gfxState.currentQuadrantMode = graphicsStateChange.fnCode
		
		case REGION_MODE_ON:
			gfxState.regionModeOn = true
			
		case REGION_MODE_OFF:
			gfxState.regionModeOn = false
			
		case END_OF_FILE:
			gfxState.fileComplete = true
			
		// For now, we're not going to do anything with any of the other ones
	}
	
	return nil
}

func (graphicsStateChange *GraphicsStateChange) String() string {
	var function string
	
	switch graphicsStateChange.fnCode {
		case REGION_MODE_ON:
			function = "Region Mode On"
		
		case REGION_MODE_OFF:
			function = "Region Mode Off"
		
		case SINGLE_QUADRANT_MODE:
			function = "Single Quadrant Mode"
		
		case MULTI_QUADRANT_MODE:
			function = "Multi Quadrant Mode"
		
		case END_OF_FILE:
			function = "End of File"
		
		case SET_UNIT_INCH:
			function = "Set Unit Inch (Warning: Deprecated)"
		
		case SET_UNIT_MM:
			function = "Set Unit MM (Warning: Deprecated)"
		
		case SET_NOTATION_ABSOLUTE:
			function = "Set Notation Absolute (Warning: Deprecated)"
		
		case SET_NOTATION_INCREMENTAL:
			function = "Set Notation Incremental (Warning: Deprecated)"
		
		case OPTIONAL_STOP:
			function = "Optional Stop (Warning: Deprecated)"
		
		case PROGRAM_STOP:
			function = "Program Stop (Warning: Deprecated)"
		
		default:
			function = "Unknown"
	}
	
	return fmt.Sprintf("{STATE CHANGE, Function: %s}", function)
}