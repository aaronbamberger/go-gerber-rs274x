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
		if move,err := interpolation.getNewCoordinate(gfxState); err != nil {
			return err
		} else {
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
							xMin := math.Min(gfxState.currentX, move.newX)
							xMax := math.Max(gfxState.currentX, move.newX)
							yMin := math.Min(gfxState.currentY, move.newY)
							yMax := math.Max(gfxState.currentY, move.newY)
							bounds.updateBounds(xMin, xMax, yMin, yMax)
							
							// Update the graphics state with the new end coordinate
							gfxState.updateCurrentCoordinate(move.newX, move.newY)
						} else {
							switch gfxState.currentInterpolationMode {
								case LINEAR_INTERPOLATION:
									// Update the bounds with both endpoints
									bounds.updateBoundsAperture(gfxState.currentX, gfxState.currentY, apertureMinSize)
									bounds.updateBoundsAperture(move.newX, move.newY, apertureMinSize)
									
									// Finally, update the graphics state with the new end coordinate
									gfxState.updateCurrentCoordinate(move.newX, move.newY)
									
								case CIRCULAR_INTERPOLATION_CLOCKWISE, CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
									radius := math.Hypot(move.newX - move.centerX, move.newY - move.centerY)
									
									// Update the bounds with both endpoints
									bounds.updateBoundsAperture(gfxState.currentX, gfxState.currentY, apertureMinSize)
									bounds.updateBoundsAperture(move.newX, move.newY, apertureMinSize)
									
									// Special case, if the angles are equal, and we're in multi quadrant mode, we're drawing a full circle,
									// so the arc spans all of the axes
									if epsilonEquals(move.startAngle, move.endAngle, gfxState.filePrecision) && (gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE) {
										bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize) // positive y-axis
										bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize) // positive x-axis
										bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize) // negative y-axis
										bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize) // negative x-axis
									} else {
										// Otherwise, if the two angles span one (or more, depending on quadrant mode) of the axes, also update the bounds with the point
										// along that axis at a distance of the radius of the arc (the max distance in that direction that the arc will cover)
										switch gfxState.currentQuadrantMode {
											case SINGLE_QUADRANT_MODE:
												if gfxState.currentInterpolationMode == CIRCULAR_INTERPOLATION_CLOCKWISE {
													if inQuadrant(move.startAngle, QUADRANT_2) && inQuadrant(move.endAngle, QUADRANT_1) {
														// The angle spans the positive y-axis
														bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
													}
													
													if inQuadrant(move.startAngle, QUADRANT_1) && inQuadrant(move.endAngle, QUADRANT_4) {
														// The angle spans the positive x-axis
														bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
													}
													
													if inQuadrant(move.startAngle, QUADRANT_4) && inQuadrant(move.endAngle, QUADRANT_3) {
														// The angle spans the negative y-axis
														bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
													}
													
													if inQuadrant(move.startAngle, QUADRANT_3) && inQuadrant(move.endAngle, QUADRANT_2) {
														// The angle spans the negative x-axis
														bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
													}
												} else {
													if inQuadrant(move.startAngle, QUADRANT_1) && inQuadrant(move.endAngle, QUADRANT_2) {
														// The angle spans the positive y-axis
														bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
													}
													
													if inQuadrant(move.startAngle, QUADRANT_4) && inQuadrant(move.endAngle, QUADRANT_1) {
														// The angle spans the positive x-axis
														bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
													}
													
													if inQuadrant(move.startAngle, QUADRANT_3) && inQuadrant(move.endAngle, QUADRANT_4) {
														// The angle spans the negative y-axis
														bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
													}
													
													if inQuadrant(move.startAngle, QUADRANT_2) && inQuadrant(move.endAngle, QUADRANT_3) {
														// The angle spans the negative x-axis
														bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
													}
												}
											
											case MULTI_QUADRANT_MODE:
												if gfxState.currentInterpolationMode == CIRCULAR_INTERPOLATION_CLOCKWISE {
													if inQuadrant(move.startAngle, QUADRANT_1) {
														if inQuadrant(move.endAngle, QUADRANT_4) {
															// The angle spans the positive x-axis
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_3) {
															// The angle spans the positive x-axis and negative y-axis
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_4) {
															// The angle spans the positive x-axis, negative y-axis, and negative x-axis
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_1) && (move.endAngle > move.startAngle) {
															// The angle spans all 4 axes
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
														}
													} else if inQuadrant(move.startAngle, QUADRANT_2) {
														if inQuadrant(move.endAngle, QUADRANT_1) {
															// The angle spans the positive y-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_4) {
															// The angle spans the positive y-axis and positive x-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_3) {
															// The angle spans the positive y-axis, positive x-axis, and negative y-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_2) && (move.endAngle > move.startAngle) {
															// The angle spans all 4 axes
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
														}
													} else if inQuadrant(move.startAngle, QUADRANT_3) {
														if inQuadrant(move.endAngle, QUADRANT_2) {
															// The angle spans the negative x-axis
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_1) {
															// The angle spans the negative x-axis and positive y-axis
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_4) {
															// The angle spans the negative x-axis, positive y-axis, and positive x-axis
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_3) && (move.endAngle > move.startAngle) {
															// The angle spans all 4 axes
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
														}
													} else if inQuadrant(move.startAngle, QUADRANT_4) {
														if inQuadrant(move.endAngle, QUADRANT_3) {
															// The angle spans the negative y-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_2) {
															// The angle spans the negative y-axis and negative x-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_1) {
															// The angle spans the negative y-axis, negative x-axis, and positive y-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_4) && (move.endAngle > move.startAngle) {
															// The angle spans all 4 axes
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
														}
													}
												} else {
													if inQuadrant(move.startAngle, QUADRANT_1) {
														if inQuadrant(move.endAngle, QUADRANT_2) {
															// The angle spans the positive y-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_3) {
															// The angle spans the positive y-axis and negative x-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_4) {
															// The angle spans the positive y-axis, negative x-axis, and negative y-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_1) && (move.endAngle < move.startAngle) {
															// The angle spans all 4 axes
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
														}
													} else if inQuadrant(move.startAngle, QUADRANT_2) {
														if inQuadrant(move.endAngle, QUADRANT_3) {
															// The angle spans the negative x-axis
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_4) {
															// The angle spans the negative x-axis and negative y-axis
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_1) {
															// The angle spans the negative x-axis, negative y-axis, and positive x-axis
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_2) && (move.endAngle < move.startAngle) {
															// The angle spans all 4 axes
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
														}
													} else if inQuadrant(move.startAngle, QUADRANT_3) {
														if inQuadrant(move.endAngle, QUADRANT_4) {
															// The angle spans the negative y-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_1) {
															// The angle spans the negative y-axis and positive x-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_2) {
															// The angle spans the negative y-axis, positive x-axis, and positive y-axis
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_3) && (move.endAngle < move.startAngle) {
															// The angle spans all 4 axes
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
														}
													} else if inQuadrant(move.startAngle, QUADRANT_4) {
														if inQuadrant(move.endAngle, QUADRANT_1) {
															// The angle spans the positive x-axis
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_2) {
															// The angle spans the positive x-axis and positive y-axis
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_3) {
															// The angle spans the positive x-axis, positive y-axis, and negative x-axis
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
														} else if inQuadrant(move.endAngle, QUADRANT_4) && (move.endAngle < move.startAngle) {
															// The angle spans all 4 axes
															bounds.updateBoundsAperture(move.centerX + radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY + radius, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX - radius, move.centerY, apertureMinSize)
															bounds.updateBoundsAperture(move.centerX, move.centerY - radius, apertureMinSize)
														}
													}
												}
										}
									}
									
									// Finally, update the graphics state with the new end coordinate
									gfxState.updateCurrentCoordinate(move.newX, move.newY)
							}
						}
					}
				
				case MOVE_OPERATION, FLASH_OPERATION:
					// For bounds checking, we treat moves and flashes the same
					if !gfxState.apertureSet {
						return fmt.Errorf("Attempt to check interpolation bounds before aperture set")
					}
					
					if aperture,found := gfxState.apertures[gfxState.currentAperture]; !found {
						return fmt.Errorf("Attempt to use aperture %d in bounds check before it has been defined", gfxState.currentAperture)
					} else {
						bounds.updateBoundsAperture(move.newX, move.newY, aperture.GetMinSize(gfxState))
						gfxState.updateCurrentCoordinate(move.newX, move.newY)	
					}
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
	if move,err := interpolation.getNewCoordinate(gfxState); err != nil {
		return err
	} else {
		switch interpolation.opCode {
			case INTERPOLATE_OPERATION:
				if !gfxState.apertureSet {
					return fmt.Errorf("Attempt to draw before aperture set")
				}
				
				if aperture,found := gfxState.apertures[gfxState.currentAperture]; !found {
					return fmt.Errorf("Attempt to use aperture %d before it has been defined", gfxState.currentAperture)
				} else {
					//apertureMinSize := aperture.GetMinSize(gfxState)
					
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
							aperture.StrokeApertureLinear(surface, gfxState, gfxState.currentX, gfxState.currentY, move.newX, move.newY)
							
							// Finally, update the graphics state with the new end coordinate
							gfxState.updateCurrentCoordinate(move.newX, move.newY)
							
						case CIRCULAR_INTERPOLATION_CLOCKWISE:
							if gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE {
								// If we're in multi-quadrant mode, we may need to adjust the end angle by 2 pi so that the angle spans come out right
								if epsilonEquals(move.startAngle, move.endAngle, gfxState.filePrecision) {
									// If the angles are equal, we're drawing a full circle, so process the end angle by 2 pi in the proper direction
									move.endAngle -= TWO_PI
								} else if move.endAngle > move.startAngle {
									// If the end angle is greater than the start angle, we process the end angle by 2 pi in the proper direction
									move.endAngle += TWO_PI
								}
							}
							/*
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
							*/
							radius := math.Hypot(move.newX - move.centerX, move.newY - move.centerY)
							aperture.StrokeApertureClockwise(surface, gfxState, move.centerX, move.centerY, radius, move.startAngle, move.endAngle)
							
							// Finally, update the graphics state with the new end coordinate
							gfxState.updateCurrentCoordinate(move.newX, move.newY)
						
						case CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
							if gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE {
								// If we're in multi-quadrant mode, we may need to adjust the end angle by 2 pi so that the angle spans come out right
								if epsilonEquals(move.startAngle, move.endAngle, gfxState.filePrecision) {
									// If the angles are equal, we're drawing a full circle, so process the end angle by 2 pi in the proper direction
									move.endAngle += TWO_PI
								} else if move.startAngle > move.endAngle {
									// If the end angle is greater than the start angle, we process the end angle by 2 pi in the proper direction
									move.endAngle += TWO_PI
								}
							}
							/*
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
							// NOTE: This is probably redundant, but because of how I'm optimizing theaperture.DrawApertureSurfaceNoHole(surface, gfxState, centerX, centerY)
							// coordinate stepping, it's possible that we won't exactly hit the end,
							// so we do it here again just in case
							if err := aperture.DrawApertureSurface(surface, gfxState, newX, newY); err != nil {
								return err
							}
							*/
							radius := math.Hypot(move.newX - move.centerX, move.newY - move.centerY)
							aperture.StrokeApertureCounterClockwise(surface, gfxState, move.centerX, move.centerY, radius, move.startAngle, move.endAngle)
							
							// Finally, update the graphics state with the new end coordinate
							gfxState.updateCurrentCoordinate(move.newX, move.newY)
					}
				}
				
			case MOVE_OPERATION:
				gfxState.updateCurrentCoordinate(move.newX, move.newY)
			
			case FLASH_OPERATION:
				if !gfxState.apertureSet {
					return fmt.Errorf("Attempt to draw before aperture set")
				}
				
				if aperture,found := gfxState.apertures[gfxState.currentAperture]; !found {
					return fmt.Errorf("Attempt to use aperture %d before it has been defined", gfxState.currentAperture)
				} else {
					gfxState.updateCurrentCoordinate(move.newX, move.newY)
					return aperture.DrawApertureSurface(surface, gfxState, gfxState.currentX, gfxState.currentY)	
				}
		}
		
		return nil
	}
}

func (interpolation *Interpolation) performDrawRegionOn(surface *cairo.Surface, gfxState *GraphicsState) error {
	if move,err := interpolation.getNewCoordinate(gfxState); err != nil {
		return err
	} else {
		switch interpolation.opCode {
			case INTERPOLATE_OPERATION:
				// Add the new segment to the current surface path
				switch gfxState.currentInterpolationMode {
					case LINEAR_INTERPOLATION:
						surface.LineTo(move.newX, move.newY)
					
					case CIRCULAR_INTERPOLATION_CLOCKWISE:
						radius := math.Hypot(gfxState.currentX - move.centerX, gfxState.currentY - move.centerY)
						if epsilonEquals(move.startAngle, move.endAngle, gfxState.filePrecision) && (gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE) {
							// NOTE: Special case, if the angles are equal, and we're in multi quadrant mode, we're drawing a full circle
							// TODO: This feels hacky, see if I can come up with a better way to handle this
							move.endAngle -= TWO_PI
						}
						
						// NOTE: The arc direction is relative to the gerber file coordinate frame
						// The conversion to the cairo coordinate frame is inherent in the y-axis mirror transformation of the surface
						surface.ArcNegative(move.centerX, move.centerY, radius, move.startAngle, move.endAngle)
					
					case CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
						radius := math.Hypot(gfxState.currentX - move.centerX, gfxState.currentY - move.centerY)
						if epsilonEquals(move.startAngle, move.endAngle, gfxState.filePrecision) && (gfxState.currentQuadrantMode == MULTI_QUADRANT_MODE) {
							// NOTE: Special case, if the angles are equal, and we're in multi quadrant mode, we're drawing a full circle
							// TODO: This feels hacky, see if I can come up with a better way to handle this
							move.endAngle += TWO_PI
						}
						
						// NOTE: The arc direction is relative to the gerber file coordinate frame
						// The conversion to the cairo coordinate frame is inherent in the y-axis mirror transformation of the surface
						surface.Arc(move.centerX, move.centerY, radius, move.startAngle, move.endAngle)
				}
			
				gfxState.updateCurrentCoordinate(move.newX, move.newY)
				
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
				gfxState.updateCurrentCoordinate(move.newX, move.newY)
			
			case FLASH_OPERATION:
				return fmt.Errorf("Flash operations are not allowed while in region mode")
		}
		
		return nil
	}
}

type InterpolationMove struct {
	newX float64
	newY float64
	centerX float64
	centerY float64
	startAngle float64
	endAngle float64
}

func (interpolation *Interpolation) getNewCoordinate(gfxState *GraphicsState) (*InterpolationMove, error) {
	newMove := new(InterpolationMove)

	// First, compute the new ending coordinates.  These are the same no matter what interpolation mode we're in
	if interpolation.xValid {
		switch gfxState.coordinateNotation {
			case ABSOLUTE_NOTATION:
				newMove.newX = interpolation.x
			
			case INCREMENTAL_NOTATION:
				newMove.newX = gfxState.currentX + interpolation.x
		}
	} else {
		newMove.newX = gfxState.currentX
	}
	
	if interpolation.yValid {
		switch gfxState.coordinateNotation {
			case ABSOLUTE_NOTATION:
				newMove.newY = interpolation.y
			
			case INCREMENTAL_NOTATION:
				newMove.newY = gfxState.currentY + interpolation.y
		}
	} else {
		newMove.newY = gfxState.currentY
	}

	switch gfxState.currentInterpolationMode {
		case LINEAR_INTERPOLATION:
			// If this is a linear interpolation, we're done (we only need to calculate center and angle for circular interpolations)
			return newMove,nil
		
		case CIRCULAR_INTERPOLATION_CLOCKWISE, CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
			switch gfxState.currentQuadrantMode {
				case SINGLE_QUADRANT_MODE:
					centerXCandidates := []float64{gfxState.currentX - interpolation.i, gfxState.currentX + interpolation.i}
					centerYCandidates := []float64{gfxState.currentY - interpolation.j, gfxState.currentY + interpolation.j}
					for _,x := range centerXCandidates {
						for _,y := range centerYCandidates {
							// First, check if the candidate center even makes sense, by checking the lengths of the line segments
							// between the candidate center and the two endpoints.  If the candidate center makes sense, these two line
							// segments will be radii of the arc, and will be equal
							radius1 := math.Hypot(gfxState.currentX - x, gfxState.currentY - y)
							radius2 := math.Hypot(newMove.newX - x, newMove.newY - y)
							if !epsilonEquals(radius1, radius2, gfxState.filePrecision) {
								continue
							}
							
							// If we have a candidate center that at least makes mathematical sense, now we need to make sure
							// it also works with the current interpolation direction
							// Start by computing the starting and ending angles of the arc
							startAngle := math.Atan2(gfxState.currentY - y, gfxState.currentX - x)
							endAngle := math.Atan2(newMove.newY - y, newMove.newX - x)
							
							// Now, make sure the the candidate center produces an arc with the correct direction that is <= 90 degrees
							// NOTE: All of the comparisons are done in the gerber-file coordinate frame
							// The conversion to the cairo coordinate frame is inherent in the y-axis mirror transformation of the surface
							switch gfxState.currentInterpolationMode {
								case CIRCULAR_INTERPOLATION_CLOCKWISE:
									if (startAngle >= endAngle) && ((startAngle - endAngle) <= ONE_HALF_PI) {
										// This covers all cases where the two angles don't straddle the -x axis (where the range of arctan changes sign)
										newMove.centerX = x
										newMove.centerY = y
										newMove.startAngle = startAngle
										newMove.endAngle = endAngle
										return newMove,nil
									} else if (startAngle < 0.0) && (endAngle > 0.0) {
										// If the two angles straddle the x-axis, the angle difference we calculate will actually be the negative complement
										// of the angle we care about, which is why the logic seems backwards (but isn't)
										if (startAngle - endAngle) <= -THREE_HALVES_PI {
											newMove.centerX = x
											newMove.centerY = y
											newMove.startAngle = startAngle
											newMove.endAngle = endAngle
											return newMove,nil
										}
									}
								
								case CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE:
									if (startAngle <= endAngle) && ((endAngle - startAngle) <= ONE_HALF_PI) {
										// This covers all cases where the two angles don't straddle the -x axis (where the range of arctan changes sign)
										newMove.centerX = x
										newMove.centerY = y
										newMove.startAngle = startAngle
										newMove.endAngle = endAngle
										return newMove,nil
									} else if (startAngle > 0.0) && (endAngle < 0.0) {
										// If the two angles straddle the x-axis, the angle difference we calculate will actually be the complement
										// of the angle we care about, which is why the logic seems backwards (but isn't)
										if (startAngle - endAngle) >= THREE_HALVES_PI {
											newMove.centerX = x
											newMove.centerY = y
											newMove.startAngle = startAngle
											newMove.endAngle = endAngle
											return newMove,nil
										}
									}
							}
						}
					}
					return nil,fmt.Errorf("Didn't find an acceptable center for single quadrant circular interpolation")
				
				case MULTI_QUADRANT_MODE:
					newMove.centerX = gfxState.currentX + interpolation.i
					newMove.centerY = gfxState.currentY + interpolation.j
					newMove.startAngle = math.Atan2(gfxState.currentY - newMove.centerY, gfxState.currentX - newMove.centerX)
					newMove.endAngle = math.Atan2(newMove.newY - newMove.centerY, newMove.newX - newMove.centerX)
					return newMove,nil
					
				default:
					return nil,fmt.Errorf("No quadrant mode set for circular interpolation")
			}
			
			default:
				return nil,fmt.Errorf("No interpolation mode set for interpolation")
	}
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