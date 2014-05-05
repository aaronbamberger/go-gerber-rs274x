package gerber_rs274x

import (
	"fmt"
	cairo "github.com/ungerik/go-cairo"
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

func (interpolation *Interpolation) ProcessDataBlockBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState) error {
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
				newX,newY := getNewCoordinate(interpolation, gfxState)
				
				if (epsilonEquals(newX, gfxState.currentX, gfxState.filePrecision)) {
					// Vertical line
					if newY > gfxState.currentY {
						for y := gfxState.currentY; y <= newY; y += gfxState.filePrecision {
							if err := drawApertureBoundsCheck(bounds, gfxState, newX, y); err != nil {
								return err
							}
						}
					} else {
						for y := gfxState.currentY; y >= newY; y -= gfxState.filePrecision {
							if err := drawApertureBoundsCheck(bounds, gfxState, newX, y); err != nil {
								return err
							}
						}
					}
				} else if (epsilonEquals(newY, gfxState.currentY, gfxState.filePrecision)) {
					// Horizontal line
					if newX > gfxState.currentX {
						for x := gfxState.currentX; x <= newX; x += gfxState.filePrecision {
							if err := drawApertureBoundsCheck(bounds, gfxState, x, newY); err != nil {
								return err
							}
						}
					} else {
						for x := gfxState.currentX; x >= newX; x -= gfxState.filePrecision {
							if err := drawApertureBoundsCheck(bounds, gfxState, x, newY); err != nil {
								return err
							}
						}
					}
				} else {
					// Any other line
				}
				
				gfxState.updateCurrentCoordinate(newX, newY)
				
			case MOVE_OPERATION:
				newX,newY := getNewCoordinate(interpolation, gfxState)
				gfxState.updateCurrentCoordinate(newX, newY)
			
			case FLASH_OPERATION:
				newX,newY := getNewCoordinate(interpolation, gfxState)
				gfxState.updateCurrentCoordinate(newX, newY)
				return drawApertureBoundsCheck(bounds, gfxState, gfxState.currentX, gfxState.currentY)
		}
	}
	
	return nil
}

func (interpolation *Interpolation) ProcessDataBlockSurface(surface *cairo.Surface, gfxState *GraphicsState) error {
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
				newX,newY := getNewCoordinate(interpolation, gfxState)
				
				if (epsilonEquals(newX, gfxState.currentX, gfxState.filePrecision)) {
					// Vertical line
					if newY > gfxState.currentY {
						for y := gfxState.currentY; y <= newY; y += gfxState.drawPrecision {
							if err := drawAperture(surface, gfxState, newX, y); err != nil {
								return err
							}
						}
					} else {
						for y := gfxState.currentY; y >= newY; y -= gfxState.drawPrecision {
							if err := drawAperture(surface, gfxState, newX, y); err != nil {
								return err
							}
						}
					}
				} else if (epsilonEquals(newY, gfxState.currentY, gfxState.filePrecision)) {
					// Horizontal line
					if newX > gfxState.currentX {
						for x := gfxState.currentX; x <= newX; x += gfxState.drawPrecision {
							if err := drawAperture(surface, gfxState, x, newY); err != nil {
								return err
							}
						}
					} else {
						for x := gfxState.currentX; x >= newX; x -= gfxState.drawPrecision {
							if err := drawAperture(surface, gfxState, x, newY); err != nil {
								return err
							}
						}
					}
				} else {
					// Any other line
				}
				
				// Make sure we draw the aperture at the actual end coordinate.
				// NOTE: This is probably redundant, but because of how I'm optimizing the
				// coordinate stepping, it's possible that we won't exactly hit the end,
				// so we do it here again just in case
				if err := drawAperture(surface, gfxState, newX, newY); err != nil {
					return err
				}
				
				gfxState.updateCurrentCoordinate(newX, newY)
				
			case MOVE_OPERATION:
				newX,newY := getNewCoordinate(interpolation, gfxState)
				gfxState.updateCurrentCoordinate(newX, newY)
			
			case FLASH_OPERATION:
				newX,newY := getNewCoordinate(interpolation, gfxState)
				gfxState.updateCurrentCoordinate(newX, newY)
				return drawAperture(surface, gfxState, gfxState.currentX, gfxState.currentY)
		}
	}
	
	return nil
}

func drawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error {
	if !gfxState.apertureSet {
		return fmt.Errorf("Attempt to use aperture before current aperture has been defined")
	}
	
	gfxState.apertures[gfxState.currentAperture].DrawApertureBoundsCheck(bounds, gfxState, x, y)

	return nil
}

func drawAperture(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	if !gfxState.apertureSet {
		return fmt.Errorf("Attempt to use aperture before current aperture has been defined")
	}
	
	gfxState.apertures[gfxState.currentAperture].DrawApertureSurface(surface, gfxState, x, y)

	return nil
}

func getNewCoordinate(interpolation *Interpolation, gfxState *GraphicsState) (newX float64, newY float64) {

	switch gfxState.currentInterpolationMode {
		case LINEAR_INTERPOLATION:
			if interpolation.xValid {
				switch gfxState.coordinateNotation {
					case ABSOLUTE_NOTATION:
						newX = interpolation.x
					
					case INCREMENTAL_NOTATION:
						newX = gfxState.currentX + interpolation.x
				}
			} else {
				newX = gfxState.currentX
			}
			
			if interpolation.yValid {
				switch gfxState.coordinateNotation {
					case ABSOLUTE_NOTATION:
						newY = interpolation.y
					
					case INCREMENTAL_NOTATION:
						newY = gfxState.currentY + interpolation.y
				}
			} else {
				newY = gfxState.currentY
			}
			
		case CIRCULAR_INTERPOLATION_CLOCKWISE:
			
		
		case CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
		
	}
	
	return
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