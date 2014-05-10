package gerber_rs274x

import (
	"fmt"
	"math"
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
				newX,newY,_,_,_,_ := interpolation.getNewCoordinate(gfxState)
			
				if gfxState.regionModeOn {
					//TODO: Do a better job than this, this is just a quick hack
					//It works for linear segments, but not for arcs
					xMin := math.Min(gfxState.currentX, newX)
					xMax := math.Max(gfxState.currentX, newX)
					yMin := math.Min(gfxState.currentY, newY)
					yMax := math.Max(gfxState.currentY, newY)
					bounds.updateBounds(xMin, xMax, yMin, yMax)
					
				} else {
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
				}
				
				gfxState.updateCurrentCoordinate(newX, newY)
				
			case MOVE_OPERATION:
				newX,newY,_,_,_,_ := interpolation.getNewCoordinate(gfxState)
				gfxState.updateCurrentCoordinate(newX, newY)
			
			case FLASH_OPERATION:
				newX,newY,_,_,_,_ := interpolation.getNewCoordinate(gfxState)
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
		if gfxState.regionModeOn {
			 return interpolation.performDrawRegionOn(surface, gfxState)
		} else {
			return interpolation.performDrawRegionOff(surface, gfxState)
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

func (interpolation *Interpolation) performDrawRegionOff(surface *cairo.Surface, gfxState *GraphicsState) error {
	switch interpolation.opCode {
		case INTERPOLATE_OPERATION:
			switch gfxState.currentInterpolationMode {
				case LINEAR_INTERPOLATION:
					newX,newY,_,_,_,_ := interpolation.getNewCoordinate(gfxState)
					lineAngle := math.Atan2(newY - gfxState.currentY, newX - gfxState.currentX)
					lineLength := math.Hypot(newX - gfxState.currentX, newY - gfxState.currentY)
					totalSteps := lineLength / gfxState.drawPrecision
					xDrawStep := gfxState.drawPrecision * math.Cos(lineAngle)
					yDrawStep := gfxState.drawPrecision * math.Sin(lineAngle)
					
					for x,y,step := gfxState.currentX,gfxState.currentY,0.0; step < totalSteps; x,y,step = x + xDrawStep,y + yDrawStep,step + 1.0 {
						if err := drawAperture(surface, gfxState, x, y); err != nil {
							return err
						}
					}
					
					// Make sure we draw the aperture at the actual end coordinate.
					// NOTE: This is probably redundant, but because of how I'm optimizing the
					// coordinate stepping, it's possible that we won't exactly hit the end,
					// so we do it here again just in case
					if err := drawAperture(surface, gfxState, newX, newY); err != nil {
						return err
					}
					
					// Finally, update the graphics state with the new end coordinate
					gfxState.updateCurrentCoordinate(newX, newY)
					
				case CIRCULAR_INTERPOLATION_CLOCKWISE:
					newX,newY,centerX,centerY,angle1,angle2 := interpolation.getNewCoordinate(gfxState)
					if epsilonEquals(angle1, angle2, gfxState.filePrecision) && (gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE) {
						// NOTE: Special case, if the angles are equal, and we're in multi quadrant mode, we're drawing a full circle
						// TODO: This feels hacky, see if I can come up with a better way to handle this
						angle2 -= (2.0 * math.Pi)
					}
					radius := math.Hypot(newX - centerX, newY - centerY)
					angleStep := gfxState.drawPrecision / radius
					
					for angle := angle1; angle > angle2; angle -= angleStep {
						offsetX := radius * math.Cos(angle)
						offsetY := radius * math.Sin(angle)
						if err := drawAperture(surface, gfxState, centerX + offsetX, centerY + offsetY); err != nil {
							return err
						} 
					}
					
					// Make sure we draw the aperture at the actual end coordinate.
					// NOTE: This is probably redundant, but because of how I'm optimizing the
					// coordinate stepping, it's possible that we won't exactly hit the end,
					// so we do it here again just in case
					if err := drawAperture(surface, gfxState, newX, newY); err != nil {
						return err
					}
					
					// Finally, update the graphics state with the new end coordinate
					gfxState.updateCurrentCoordinate(newX, newY)
				
				case CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
					newX,newY,centerX,centerY,angle1,angle2 := interpolation.getNewCoordinate(gfxState)
					if epsilonEquals(angle1, angle2, gfxState.filePrecision) && (gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE) {
						// NOTE: Special case, if the angles are equal, and we're in multi quadrant mode, we're drawing a full circle
						// TODO: This feels hacky, see if I can come up with a better way to handle this
						angle2 += (2.0 * math.Pi)
					}
					radius := math.Hypot(newX - centerX, newY - centerY)
					angleStep := gfxState.drawPrecision / radius
					
					for angle := angle1; angle < angle2; angle += angleStep {
						offsetX := radius * math.Cos(angle)
						offsetY := radius * math.Sin(angle)
						if err := drawAperture(surface, gfxState, centerX + offsetX, centerY + offsetY); err != nil {
							return err
						} 
					}
					
					// Make sure we draw the aperture at the actual end coordinate.
					// NOTE: This is probably redundant, but because of how I'm optimizing the
					// coordinate stepping, it's possible that we won't exactly hit the end,
					// so we do it here again just in case
					if err := drawAperture(surface, gfxState, newX, newY); err != nil {
						return err
					}
					
					// Finally, update the graphics state with the new end coordinate
					gfxState.updateCurrentCoordinate(newX, newY)
			}
			
		case MOVE_OPERATION:
			newX,newY,_,_,_,_ := interpolation.getNewCoordinate(gfxState)
			gfxState.updateCurrentCoordinate(newX, newY)
		
		case FLASH_OPERATION:
			newX,newY,_,_,_,_ := interpolation.getNewCoordinate(gfxState)
			gfxState.updateCurrentCoordinate(newX, newY)
			return drawAperture(surface, gfxState, gfxState.currentX, gfxState.currentY)
	}
	
	return nil
}

func (interpolation *Interpolation) performDrawRegionOn(surface *cairo.Surface, gfxState *GraphicsState) error {
	switch interpolation.opCode {
		case INTERPOLATE_OPERATION:
			newX,newY,centerX,centerY,angle1,angle2 := interpolation.getNewCoordinate(gfxState)
			
			// Add the new segment to the current surface path
			switch gfxState.currentInterpolationMode {
				case LINEAR_INTERPOLATION:
					correctedX := (newX * gfxState.scaleFactor) + gfxState.xOffset
					correctedY := (newY * gfxState.scaleFactor) + gfxState.yOffset
					surface.LineTo(correctedX, correctedY)
				
				case CIRCULAR_INTERPOLATION_CLOCKWISE:
					correctedCenterX := (centerX * gfxState.scaleFactor) + gfxState.xOffset
					correctedCenterY := (centerY * gfxState.scaleFactor) + gfxState.yOffset
					scaledRadius := math.Hypot(gfxState.currentX - centerX, gfxState.currentY - centerY) * gfxState.scaleFactor
					if epsilonEquals(angle1, angle2, gfxState.filePrecision) && (gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE) {
						// NOTE: Special case, if the angles are equal, and we're in multi quadrant mode, we're drawing a full circle
						// TODO: This feels hacky, see if I can come up with a better way to handle this
						angle2 -= (2.0 * math.Pi)
					}
					
					// NOTE: The arc direction is relative to the gerber file coordinate frame
					// The conversion to the cairo coordinate frame is inherent in the y-axis mirror transformation of the surface
					surface.ArcNegative(correctedCenterX, correctedCenterY, scaledRadius, angle1, angle2)
				
				case CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
					correctedCenterX := (centerX * gfxState.scaleFactor) + gfxState.xOffset
					correctedCenterY := (centerY * gfxState.scaleFactor) + gfxState.yOffset
					scaledRadius := math.Hypot(gfxState.currentX - centerX, gfxState.currentY - centerY) * gfxState.scaleFactor
					if epsilonEquals(angle1, angle2, gfxState.filePrecision) && (gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE) {
						// NOTE: Special case, if the angles are equal, and we're in multi quadrant mode, we're drawing a full circle
						// TODO: This feels hacky, see if I can come up with a better way to handle this
						angle2 += (2.0 * math.Pi)
					}
					
					// NOTE: The arc direction is relative to the gerber file coordinate frame
					// The conversion to the cairo coordinate frame is inherent in the y-axis mirror transformation of the surface
					surface.Arc(correctedCenterX, correctedCenterY, scaledRadius, angle1, angle2)
			}
		
			gfxState.updateCurrentCoordinate(newX, newY)
			
		case MOVE_OPERATION:
			// If we're in region mode, this means we're closing off a contour.  First, set the proper polarity,
			// then perform the actual draw
			switch gfxState.currentLevelPolarity {
				case DARK_POLARITY:
					surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
				
				case CLEAR_POLARITY:
					surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
			}
			surface.Fill()
			
			// Now, update the current point
			newX,newY,_,_,_,_ := interpolation.getNewCoordinate(gfxState)
			gfxState.updateCurrentCoordinate(newX, newY)
		
		case FLASH_OPERATION:
			return fmt.Errorf("Flash operations are not allowed while in region mode")
	}
	
	return nil
}

func (interpolation *Interpolation) getNewCoordinate(gfxState *GraphicsState) (newX float64, newY float64, centerX float64, centerY float64, angle1 float64, angle2 float64) {

	// First, compute the new ending coordinates.  These are the same no matter what interpolation mode we're in
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

	switch gfxState.currentInterpolationMode {
		case LINEAR_INTERPOLATION:
			// If this is a linear interpolation, we're done (we only need to calculate center and angle for circular interpolations)
			return newX,newY,0.0,0.0,0.0,0.0
		
		case CIRCULAR_INTERPOLATION_CLOCKWISE, CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
			switch gfxState.currentQuadrantMode {
				case SINGLE_QUADRANT_MODE:
					centerXCandidates := []float64{gfxState.currentX - interpolation.i, gfxState.currentX + interpolation.i}
					centerYCandidates := []float64{gfxState.currentY - interpolation.j, gfxState.currentY + interpolation.j}
					for _,x := range centerXCandidates {
						for _,y := range centerYCandidates {
							// Compute the starting and ending angles of the arc, and then the circumscribed angle
							angle1 := math.Atan2(gfxState.currentY - y, gfxState.currentX - x)
							angle2 := math.Atan2(newY - y, newX - x)
							theta := angle2 - angle1
							
							if math.Abs(theta) <= (math.Pi / 2.0) {
								// If the angle is <= 90 degrees, check whether it's the correct direction
								// for the current interpolation mode.  If it is, we've found the correct center, so return
								// NOTE: All of the comparisons are done in the gerber-file coordinate frame
								// The conversion to the cairo coordinate frame is inherent in the y-axis mirror transformation of the surface
								if gfxState.currentInterpolationMode == CIRCULAR_INTERPOLATION_CLOCKWISE {
									if angle1 > angle2 {
										return newX,newY,x,y,angle1,angle2
									}
								} else {
									if angle2 > angle1 {
										return newX,newY,x,y,angle1,angle2
									}
								}
							}
						}
					}
					fmt.Printf("ERROR: Didn't find an acceptable center\n")
				
				case MULTI_QUADRANT_MODE:
					arcCenterX := gfxState.currentX + interpolation.i
					arcCenterY := gfxState.currentY + interpolation.j
					angle1 := math.Atan2(gfxState.currentY - arcCenterY, gfxState.currentX - arcCenterX)
					angle2 := math.Atan2(newY - arcCenterY, newX - arcCenterX)
					return newX,newY,arcCenterX,arcCenterY,angle1,angle2
			}
	}
	
	return
}

func lawOfCosines(aX float64, aY float64, bX float64, bY float64, cX float64, cY float64) (angle float64) {
	// Use the law of cosines to compute an interior angle of a triangle, given all 3 points
	sideA := math.Hypot(bX - cX, bY - cY)
	sideB := math.Hypot(aX - cX, aY - cY)
	sideC := math.Hypot(aX - bX, aY - bY)
	
	return math.Acos((math.Pow(sideA, 2) + math.Pow(sideB, 2) - math.Pow(sideC, 2)) / (2 * sideA * sideB))
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