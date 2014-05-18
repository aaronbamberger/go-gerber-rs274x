package gerber_rs274x

import (
	"fmt"
	"math"
	cairo "github.com/ungerik/go-cairo"
)

type CircleAperture struct {
	apertureNumber int
	diameter float64
	Hole
}

func (aperture *CircleAperture) AperturePlaceholder() {

}

func (aperture *CircleAperture) GetApertureNumber() int {
	return aperture.apertureNumber
}

func (aperture *CircleAperture) GetHole() Hole {
	return aperture.Hole
}

func (aperture *CircleAperture) SetHole(hole Hole) {
	aperture.Hole = hole
}

func (aperture *CircleAperture) GetMinSize(gfxState *GraphicsState) float64 {
	return aperture.diameter / 2.0
}

func (aperture *CircleAperture) DrawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error {
	radius := aperture.diameter / 2.0
	xMin := x - radius
	xMax := x + radius
	yMin := y - radius
	yMax := y + radius
	
	bounds.updateBounds(xMin, xMax, yMin, yMax)
	
	return nil
}

func (aperture *CircleAperture) DrawApertureSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {

	radius := aperture.diameter / 2.0
	correctedX := (x - radius) * gfxState.scaleFactor
	correctedY := (y - radius) * gfxState.scaleFactor
	
	return renderApertureToSurface(aperture, surface, gfxState, correctedX, correctedY)
}

func (aperture *CircleAperture) StrokeApertureLinear(surface *cairo.Surface, gfxState *GraphicsState, startX float64, startY float64, endX float64, endY float64) error {
	
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}
	
	radius := aperture.diameter / 2.0
	strokeLength := math.Hypot(endX - startX, endY - startY)
	strokeAngle := math.Atan2(endY - startY, endX - startX)
	
	if aperture.Hole != nil && strokeLength < radius{
		// If this aperture has a hole, and the distance between the start and end of the stroke is less than the aperture radius,
		// we can't use our optimized draw because the hole won't be completely covered up in the middle of the stroke, so we fall back
		// to manually stroking the aperture
		totalSteps := strokeLength / gfxState.filePrecision
		xDrawStep := gfxState.filePrecision * math.Cos(strokeAngle)
		yDrawStep := gfxState.filePrecision * math.Sin(strokeAngle)
		
		for x,y,step := startX,startY,0.0; step < totalSteps; x,y,step = x + xDrawStep,y + yDrawStep,step + 1.0 {
			if err := aperture.DrawApertureSurface(surface, gfxState, x, y); err != nil {
				return err
			}
		}
	} else {
		// Else, we can optimize by drawing a line the thickness of the aperture diameter between the two points, then flashing the
		// aperture at each end to get the endcaps correct
		topAngle := strokeAngle + (math.Pi / 2.0)
		bottomAngle := strokeAngle - (math.Pi / 2.0)
		topOffsetX := radius * math.Cos(topAngle)
		topOffsetY := radius * math.Sin(topAngle)
		bottomOffsetX := radius * math.Cos(bottomAngle)
		bottomOffsetY := radius * math.Sin(bottomAngle)
		
		topLeftX := (startX + topOffsetX) * gfxState.scaleFactor
		topLeftY := (startY + topOffsetY) * gfxState.scaleFactor
		topRightX := (endX + topOffsetX) * gfxState.scaleFactor
		topRightY := (endY + topOffsetY) * gfxState.scaleFactor
		bottomLeftX := (startX + bottomOffsetX) * gfxState.scaleFactor
		bottomLeftY := (startY + bottomOffsetY) * gfxState.scaleFactor
		bottomRightX := (endX + bottomOffsetX) * gfxState.scaleFactor
		bottomRightY := (endY + bottomOffsetY) * gfxState.scaleFactor
		
		// Draw the stroke, except for the endpoints
		surface.MoveTo(topLeftX, topLeftY)
		surface.LineTo(topRightX, topRightY)
		surface.LineTo(bottomRightX, bottomRightY)
		surface.LineTo(bottomLeftX, bottomLeftY)
		surface.LineTo(topLeftX, topLeftY)
		surface.Fill()
		
		// Draw each of the endpoints by flashing the aperture at the endpoints
		aperture.DrawApertureSurface(surface, gfxState, startX, startY)
		aperture.DrawApertureSurface(surface, gfxState, endX, endY)
	}

	return nil
}

func (aperture *CircleAperture) StrokeApertureClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *CircleAperture) StrokeApertureCounterClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *CircleAperture) renderApertureToGraphicsState(gfxState *GraphicsState) {
	// This will render the aperture to a cairo surface the first time it is needed, then
	// cache it in the graphics state.  Subsequent draws of the aperture will used the cached surface
	
	scaledDiameter := aperture.diameter * gfxState.scaleFactor
	
	// Construct the surface we're drawing to
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, int(math.Ceil(scaledDiameter)), int(math.Ceil(scaledDiameter)))
	surface.SetAntialias(cairo.ANTIALIAS_GRAY)
	
	// Draw the aperture
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}
	
	centerX := scaledDiameter / 2.0
	centerY := scaledDiameter / 2.0
	radius := scaledDiameter / 2.0
	
	surface.Arc(centerX, centerY, radius, 0, 2.0 * math.Pi)
	surface.Fill()
	
	// If present, remove the hole
	if aperture.Hole != nil {
		aperture.DrawHoleSurface(surface, gfxState, centerX, centerY)
	}
	
	surface.WriteToPNG(fmt.Sprintf("Aperture-%d.png", aperture.apertureNumber))
	
	gfxState.renderedApertures[aperture.apertureNumber] = surface
}

func (aperture *CircleAperture) String() string {
	return fmt.Sprintf("{CA, Diameter: %f, Hole: %v}", aperture.diameter, aperture.Hole)
}