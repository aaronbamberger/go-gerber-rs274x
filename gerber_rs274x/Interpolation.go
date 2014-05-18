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
		newX,newY,centerX,centerY,angle1,angle2 := interpolation.getNewCoordinate(gfxState)
		
		switch interpolation.opCode {
			case INTERPOLATE_OPERATION:
				if !gfxState.apertureSet {
					return fmt.Errorf("Attempt to check interpolation bounds before aperture set")
				}
				
				if aperture,found := gfxState.apertures[gfxState.currentAperture]; !found {
					return fmt.Errorf("Attempt to use aperture %d in bounds check before it has been defined", gfxState.currentAperture)
				} else {
					apertureMinSize := aperture.GetMinSize(gfxState)
					
					if gfxState.regionModeOn {
						//TODO: Do a better job than this, this is just a quick hack
						//It works for linear segments, but not for arcs
						xMin := math.Min(gfxState.currentX, newX)
						xMax := math.Max(gfxState.currentX, newX)
						yMin := math.Min(gfxState.currentY, newY)
						yMax := math.Max(gfxState.currentY, newY)
						bounds.updateBounds(xMin, xMax, yMin, yMax)
						
					} else {
						switch gfxState.currentInterpolationMode {
							case LINEAR_INTERPOLATION:
								// Update the bounds with the new endpoint (the start point was accounted for in a previous operation)
								bounds.updateBoundsAperture(newX, newY, apertureMinSize)
								
								// Finally, update the graphics state with the new end coordinate
								gfxState.updateCurrentCoordinate(newX, newY)
								
							case CIRCULAR_INTERPOLATION_CLOCKWISE, CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
								radius := math.Hypot(newX - centerX, newY - centerY)
								point1X := centerX + (math.Cos(angle1) * radius)
								point1Y := centerY + (math.Sin(angle1) * radius)
								point2X := centerX + (math.Cos(angle2) * radius)
								point2Y := centerY + (math.Sin(angle2) * radius)
								
								// Update the bounds with both endpoints
								bounds.updateBoundsAperture(point1X, point1Y, apertureMinSize)
								bounds.updateBoundsAperture(point2X, point2Y, apertureMinSize)
								
								// Special case, if the angles are equal, and we're in multi quadrant mode, we're drawing a full circle,
								// so the arc spans all of the axes
								if epsilonEquals(angle1, angle2, gfxState.filePrecision) && (gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE) {
									bounds.updateBoundsAperture(centerX, centerY + radius, apertureMinSize) // positive y-axis
									bounds.updateBoundsAperture(centerX + radius, centerY, apertureMinSize) // positive x-axis
									bounds.updateBoundsAperture(centerX, centerY - radius, apertureMinSize) // negative y-axis
									bounds.updateBoundsAperture(centerX - radius, centerY, apertureMinSize) // negative x-axis
								} else {
									// Otherwise, if the two angles span one (or more, depending on quadrant mode) of the axes, also update the bounds with the point
									// along that axis at a distance of the radius of the arc (the max distance in that direction that the arc will cover)
									if gfxState.currentInterpolationMode == CIRCULAR_INTERPOLATION_CLOCKWISE {
										if (angle1 >= (math.Pi / 2.0)) && (angle2 <= (math.Pi / 2.0)) {
											// The angle spans the positive y-axis
											bounds.updateBoundsAperture(centerX, centerY + radius, apertureMinSize)
										}
										
										if (angle1 >= 0.0) && (angle2 <= 0.0) {
											// The angle spans the positive x-axis
											bounds.updateBoundsAperture(centerX + radius, centerY, apertureMinSize)
										}
										
										if (angle1 >= -(math.Pi / 2.0)) && (angle2 <= -(math.Pi / 2.0)) {
											// The angle spans the negative y-axis
											bounds.updateBoundsAperture(centerX, centerY - radius, apertureMinSize)
										}
										
										if (angle1 > -math.Pi) && (angle2 < math.Pi) {
											// The angle spans the negative x-axis
											bounds.updateBoundsAperture(centerX - radius, centerY, apertureMinSize)
										}
									} else {
										if (angle1 <= (math.Pi / 2.0)) && (angle2 >= (math.Pi / 2.0)) {
											// The angle spans the positive y-axis
											bounds.updateBoundsAperture(centerX, centerY + radius, apertureMinSize)
										}
										
										if (angle1 <= 0.0) && (angle2 >= 0.0) {
											// The angle spans the positive x-axis
											bounds.updateBoundsAperture(centerX + radius, centerY, apertureMinSize)
										}
										
										if (angle1 <= -(math.Pi / 2.0)) && (angle2 >= -(math.Pi / 2.0)) {
											// The angle spans the negative y-axis
											bounds.updateBoundsAperture(centerX, centerY - radius, apertureMinSize)
										}
										
										if (angle1 < math.Pi) && (angle2 > -math.Pi) {
											// The angle spans the negative x-axis
											bounds.updateBoundsAperture(centerX - radius, centerY, apertureMinSize)
										}
									}
								}
								
								// Finally, update the graphics state with the new end coordinate
								gfxState.updateCurrentCoordinate(newX, newY)
						}
					}
				}
				
			case MOVE_OPERATION:
				// Since this is just a move, the mins and maxes are the same
				bounds.updateBounds(newX, newX, newY, newY)
				gfxState.updateCurrentCoordinate(newX, newY)
			
			case FLASH_OPERATION:
				if !gfxState.apertureSet {
					return fmt.Errorf("Attempt to check interpolation bounds before aperture set")
				}
				
				if aperture,found := gfxState.apertures[gfxState.currentAperture]; !found {
					return fmt.Errorf("Attempt to use aperture %d in bounds check before it has been defined", gfxState.currentAperture)
				} else {
					bounds.updateBoundsAperture(newX, newY, aperture.GetMinSize(gfxState))
					gfxState.updateCurrentCoordinate(newX, newY)	
				}
		}
	}
	
	return nil
}

func (interpolation *Interpolation) ProcessDataBlockSurface(surface *cairo.Surface, gfxState *GraphicsState) error {
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

func (interpolation *Interpolation) performDrawRegionOff(surface *cairo.Surface, gfxState *GraphicsState) error {
	newX,newY,centerX,centerY,angle1,angle2 := interpolation.getNewCoordinate(gfxState)

	switch interpolation.opCode {
		case INTERPOLATE_OPERATION:
			if !gfxState.apertureSet {
				return fmt.Errorf("Attempt to draw before aperture set")
			}
			
			if aperture,found := gfxState.apertures[gfxState.currentAperture]; !found {
				return fmt.Errorf("Attempt to use aperture %d before it has been defined", gfxState.currentAperture)
			} else {
				apertureMinSize := aperture.GetMinSize(gfxState)
				
				switch gfxState.currentInterpolationMode {
					case LINEAR_INTERPOLATION:
						/*
						lineAngle := math.Atan2(newY - gfxState.currentY, newX - gfxState.currentX)
						lineLength := math.Hypot(newX - gfxState.currentX, newY - gfxState.currentY)
						totalSteps := lineLength / apertureMinSize
						xDrawStep := apertureMinSize * math.Cos(lineAngle)
						yDrawStep := apertureMinSize * math.Sin(lineAngle)
						
						for x,y,step := gfxState.currentX,gfxState.currentY,0.0; step < totalSteps; x,y,step = x + xDrawStep,y + yDrawStep,step + 1.0 {
							if err := aperture.DrawApertureSurface(surface, gfxState, x, y); err != nil {
								return err
							}
						}
						
						// Make sure we draw the aperture at the actual end coordinate.
						// NOTE: This is probably redundant, but because of how I'm optimizing the
						// coordinate stepping, it's possible that we won't exactly hit the end,
						// so we do it here again just in case
						if err := aperture.DrawApertureSurface(surface, gfxState, newX, newY); err != nil {
							return err
						}
						*/
						aperture.StrokeApertureLinear(surface, gfxState, gfxState.currentX, gfxState.currentY, newX, newY)
						
						// Finally, update the graphics state with the new end coordinate
						gfxState.updateCurrentCoordinate(newX, newY)
						
					case CIRCULAR_INTERPOLATION_CLOCKWISE:
						if epsilonEquals(angle1, angle2, gfxState.filePrecision) && (gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE) {
							// NOTE: Special case, if the angles are equal, and we're in multi quadrant mode, we're drawing a full circle
							// TODO: This feels hacky, see if I can come up with a better way to handle this
							angle2 -= (2.0 * math.Pi)
						}
						radius := math.Hypot(newX - centerX, newY - centerY)
						angleStep := apertureMinSize / radius
						
						for angle := angle1; angle > angle2; angle -= angleStep {
							offsetX := radius * math.Cos(angle)
							offsetY := radius * math.Sin(angle)
							if err := aperture.DrawApertureSurface(surface, gfxState, centerX + offsetX, centerY + offsetY); err != nil {
								return err
							}
						}
						
						// Make sure we draw the aperture at the actual end coordinate.
						// NOTE: This is probably redundant, but because of how I'm optimizing the
						// coordinate stepping, it's possible that we won't exactly hit the end,
						// so we do it here again just in case
						if err := aperture.DrawApertureSurface(surface, gfxState, newX, newY); err != nil {
							return err
						}
						
						// Finally, update the graphics state with the new end coordinate
						gfxState.updateCurrentCoordinate(newX, newY)
					
					case CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
						if epsilonEquals(angle1, angle2, gfxState.filePrecision) && (gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE) {
							// NOTE: Special case, if the angles are equal, and we're in multi quadrant mode, we're drawing a full circle
							// TODO: This feels hacky, see if I can come up with a better way to handle this
							angle2 += (2.0 * math.Pi)
						}
						radius := math.Hypot(newX - centerX, newY - centerY)
						angleStep := apertureMinSize / radius
						
						for angle := angle1; angle < angle2; angle += angleStep {
							offsetX := radius * math.Cos(angle)
							offsetY := radius * math.Sin(angle)
							if err := aperture.DrawApertureSurface(surface, gfxState, centerX + offsetX, centerY + offsetY); err != nil {
								return err
							}
						}
						
						// Make sure we draw the aperture at the actual end coordinate.
						// NOTE: This is probably redundant, but because of how I'm optimizing the
						// coordinate stepping, it's possible that we won't exactly hit the end,
						// so we do it here again just in case
						if err := aperture.DrawApertureSurface(surface, gfxState, newX, newY); err != nil {
							return err
						}
						
						// Finally, update the graphics state with the new end coordinate
						gfxState.updateCurrentCoordinate(newX, newY)
				}
			}
			
		case MOVE_OPERATION:
			gfxState.updateCurrentCoordinate(newX, newY)
		
		case FLASH_OPERATION:
			if !gfxState.apertureSet {
				return fmt.Errorf("Attempt to draw before aperture set")
			}
			
			if aperture,found := gfxState.apertures[gfxState.currentAperture]; !found {
				return fmt.Errorf("Attempt to use aperture %d before it has been defined", gfxState.currentAperture)
			} else {
				gfxState.updateCurrentCoordinate(newX, newY)
				return aperture.DrawApertureSurface(surface, gfxState, gfxState.currentX, gfxState.currentY)	
			}
	}
	
	return nil
}

func (interpolation *Interpolation) performDrawRegionOn(surface *cairo.Surface, gfxState *GraphicsState) error {
	newX,newY,centerX,centerY,angle1,angle2 := interpolation.getNewCoordinate(gfxState)

	switch interpolation.opCode {
		case INTERPOLATE_OPERATION:
			// Add the new segment to the current surface path
			switch gfxState.currentInterpolationMode {
				case LINEAR_INTERPOLATION:
					correctedX := newX * gfxState.scaleFactor
					correctedY := newY * gfxState.scaleFactor
					surface.LineTo(correctedX, correctedY)
				
				case CIRCULAR_INTERPOLATION_CLOCKWISE:
					correctedCenterX := centerX * gfxState.scaleFactor
					correctedCenterY := centerY * gfxState.scaleFactor
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
					correctedCenterX := centerX * gfxState.scaleFactor
					correctedCenterY := centerY * gfxState.scaleFactor
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