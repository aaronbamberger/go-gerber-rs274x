package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

type Interpolation struct {
	fnCode FunctionCode
	opCode OperationCode
	x float64
	y float64
	i float64
	j float64
	fnCodeValid bool
	opCodeValid bool
	xValid bool
	yValid bool
}

func (interpolation *Interpolation) DataBlockPlaceholder() {

}

func (interpolation *Interpolation) ProcessDataBlockSVG(svg *svg.SVG, gfxState *GraphicsState) error {
	//TODO: This is the hard one
	
	// First, if this interpolation has a valid function code, update the graphics state
	if interpolation.fnCodeValid {
		switch interpolation.fnCode {
			case LINEAR_INTERPOLATION, CIRCULAR_INTERPOLATION_CLOCKWISE, CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
				gfxState.currentInterpolationMode = interpolation.fnCode
				gfxState.interpolationModeSet = true
		}
	}
	
	// Next, if this interpolation has a valid operation code, perform the operation
	if interpolation.opCodeValid {
		switch interpolation.opCode {
			case INTERPOLATE_OPERATION:
				
			
			case MOVE_OPERATION:
				updateCurrentCoordinate(interpolation, gfxState)
			
			case FLASH_OPERATION:
				updateCurrentCoordinate(interpolation, gfxState)
			
				if gfxState.apertureSet {
					return drawAperture(svg, gfxState.apertures[gfxState.currentAperture], gfxState.currentX, gfxState.currentY)
				} else {
					return fmt.Errorf("Attempt to flash aperture before current aperture has been defined")
				}
		}
	}
	
	return nil
}

func updateCurrentCoordinate(interpolation *Interpolation, gfxState *GraphicsState) {
	if interpolation.xValid {
		switch gfxState.coordinateNotation {
			case ABSOLUTE_NOTATION:
				gfxState.currentX = interpolation.x
			
			case INCREMENTAL_NOTATION:
				gfxState.currentX += interpolation.x
		}
	}
	
	if interpolation.yValid {
		switch gfxState.coordinateNotation {
			case ABSOLUTE_NOTATION:
				gfxState.currentY = interpolation.y
			
			case INCREMENTAL_NOTATION:
				gfxState.currentY += interpolation.y
		}
	}
}

func drawAperture(svg *svg.SVG, aperture Aperture, x float64, y float64) error {

	return nil
}

func (interpolation *Interpolation) String() string {
	var function string
	var operation string
	
	if interpolation.fnCodeValid {
		switch interpolation.fnCode {
			case LINEAR_INTERPOLATION:
				function = "Linear Interpolation"
			
			case CIRCULAR_INTERPOLATION_CLOCKWISE:
				function = "Circular CW Interpolation"
			
			case CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
				function = "Circular CCW Interpolation"
				
			case SELECT_APERTURE:
				function = "Select Aperture (Warning: Deprecated)"
			
			case PREPARE_FOR_FLASH:
				function = "Prepare for Flash (Warning: Deprecated)"
			
			default:
				function = "Unknown"
		}
	} else {
		function = "None"
	}
	
	if interpolation.opCodeValid {
		switch interpolation.opCode {
			case INTERPOLATE_OPERATION:
				operation = "Interpolate"
			
			case MOVE_OPERATION:
				operation = "Move"
			
			case FLASH_OPERATION:
				operation = "Flash"
			
			default:
				operation = "Unknown"
		}
	} else {
		operation = "None"
	}
	
	return fmt.Sprintf("{INTERPOLATION, Fn: %s, Op: %s, X Valid: %t, X: %f, Y Valid: %t, Y: %f, I: %f, J: %f}",
						function,
						operation,
						interpolation.xValid,
						interpolation.x,
						interpolation.yValid,
						interpolation.y,
						interpolation.i,
						interpolation.j)
}