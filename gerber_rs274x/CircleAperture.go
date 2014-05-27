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
	correctedX := x - radius
	correctedY := y - radius
	
	return renderApertureToSurface(aperture, surface, gfxState, correctedX, correctedY)
}

func (aperture *CircleAperture) DrawApertureSurfaceNoHole(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	radius := aperture.diameter / 2.0
	correctedX := x - radius
	correctedY := y - radius
	
	return renderApertureNoHoleToSurface(aperture, surface, gfxState, correctedX, correctedY)
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
	
	if aperture.Hole != nil && strokeLength < radius {
		// If this aperture has a hole, and the distance between the start and end of the stroke is less than the aperture radius,
		// we can't use our optimized draw because the hole won't be completely covered up in the middle of the stroke, so we fall back
		// to manually stroking the aperture
		drawStep := strokeLength / float64(SLOW_DRAWING_STEPS)
		xDrawStep := drawStep * math.Cos(strokeAngle)
		yDrawStep := drawStep * math.Sin(strokeAngle)
		
		for x,y,step := startX,startY,0; step < SLOW_DRAWING_STEPS; x,y,step = x + xDrawStep,y + yDrawStep,step + 1 {
			if err := aperture.DrawApertureSurface(surface, gfxState, x, y); err != nil {
				return err
			}
		}
	} else {
		// Else, we can optimize by drawing a line the thickness of the aperture diameter between the two points, then flashing the
		// aperture at each end to get the endcaps correct
		topAngle := strokeAngle + ONE_HALF_PI
		bottomAngle := strokeAngle - ONE_HALF_PI
		topOffsetX := radius * math.Cos(topAngle)
		topOffsetY := radius * math.Sin(topAngle)
		bottomOffsetX := radius * math.Cos(bottomAngle)
		bottomOffsetY := radius * math.Sin(bottomAngle)
		
		topLeftX := startX + topOffsetX
		topLeftY := startY + topOffsetY
		topRightX := endX + topOffsetX
		topRightY := endY + topOffsetY
		bottomLeftX := startX + bottomOffsetX
		bottomLeftY := startY + bottomOffsetY
		bottomRightX := endX + bottomOffsetX
		bottomRightY := endY + bottomOffsetY
		
		// Draw the stroke, except for the endpoints
		surface.MoveTo(topLeftX, topLeftY)
		surface.LineTo(topRightX, topRightY)
		surface.LineTo(bottomRightX, bottomRightY)
		surface.LineTo(bottomLeftX, bottomLeftY)
		surface.LineTo(topLeftX, topLeftY)
		surface.Fill()
		
		// Draw each of the endpoints by flashing the aperture at the endpoints
		// We use the special "no hole" version of the draw, because any holes will
		// have been covered over by the rest of the aperture during the stroke
		aperture.DrawApertureSurfaceNoHole(surface, gfxState, startX, startY)
		aperture.DrawApertureSurfaceNoHole(surface, gfxState, endX, endY)
	}

	return nil
}

func (aperture *CircleAperture) StrokeApertureClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	//TODO: For testing, makes it look better for now
	surface.SetAntialias(cairo.ANTIALIAS_DEFAULT)

	strokeLength := math.Abs(startAngle - endAngle) * radius
	apertureRadius := aperture.diameter / 2.0
	
	if aperture.Hole != nil && strokeLength < apertureRadius {
		angleStep := (strokeLength / float64(SLOW_DRAWING_STEPS)) / radius
		// If this aperture has a hole, and the distance between the start and end of the stroke is less than the aperture radius,
		// we can't use our optimized draw because the hole won't be completely covered up in the middle of the stroke, so we fall back
		// to manually stroking the aperture				
		for angle := startAngle; angle > endAngle; angle -= angleStep {
			offsetX := radius * math.Cos(angle)
			offsetY := radius * math.Sin(angle)
			if err := aperture.DrawApertureSurface(surface, gfxState, centerX + offsetX, centerY + offsetY); err != nil {
				return err
			}
		}
	} else {
		// Else, we can optimize by drawing an arc the thickness of the aperture diameter between the two points, then flashing the
		// aperture at each end to get the endcaps correct
		
		// Draw the stroke, except for the endpoints	
		outerRadius := radius + apertureRadius
		innerRadius := radius - apertureRadius
		arc1StartPointX := centerX + (outerRadius * math.Cos(startAngle))
		arc1StartPointY := centerY + (outerRadius * math.Sin(startAngle))
		arc2StartPointX := centerX + (innerRadius * math.Cos(endAngle))
		arc2StartPointY := centerY + (innerRadius * math.Sin(endAngle))
		surface.MoveTo(arc1StartPointX, arc1StartPointY)
		surface.ArcNegative(centerX, centerY, outerRadius, startAngle, endAngle)
		surface.LineTo(arc2StartPointX, arc2StartPointY)
		surface.Arc(centerX, centerY, innerRadius, endAngle, startAngle)
		surface.LineTo(arc1StartPointX, arc1StartPointY)
		surface.Fill()
		
		// Draw each of the endpoints by flashing the aperture at the endpoints
		startX := centerX + (radius * math.Cos(startAngle))
		startY := centerY + (radius * math.Sin(startAngle))
		endX := centerX + (radius * math.Cos(endAngle))
		endY := centerY + (radius * math.Sin(endAngle))
		// We use the special "no hole" version of the draw, because any holes will
		// have been covered over by the rest of the aperture during the stroke
		aperture.DrawApertureSurfaceNoHole(surface, gfxState, startX, startY)
		aperture.DrawApertureSurfaceNoHole(surface, gfxState, endX, endY)
		
		fmt.Printf("Center (%f %f), Start (%f %f), End (%f %f)\n", centerX, centerY, startX, startY, endX, endY)
	}
	
	//TODO: Reset so other draw operations can make their own antialiasing decisions
	surface.SetAntialias(cairo.ANTIALIAS_NONE)

	return nil
}

func (aperture *CircleAperture) StrokeApertureCounterClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	//TODO: For testing, makes it look better for now
	surface.SetAntialias(cairo.ANTIALIAS_DEFAULT)

	strokeLength := math.Abs(startAngle - endAngle) * radius
	apertureRadius := aperture.diameter / 2.0
	
	fmt.Printf("Start angle %f, End angle %f, Stroke Length %f, Aperture Radius %f\n", startAngle, endAngle, strokeLength, apertureRadius)
	
	if aperture.Hole != nil && strokeLength < apertureRadius {
		angleStep := (strokeLength / float64(SLOW_DRAWING_STEPS)) / radius
		// If this aperture has a hole, and the distance between the start and end of the stroke is less than the aperture radius,
		// we can't use our optimized draw because the hole won't be completely covered up in the middle of the stroke, so we fall back
		// to manually stroking the aperture				
		for angle := startAngle; angle < endAngle; angle += angleStep {
			offsetX := radius * math.Cos(angle)
			offsetY := radius * math.Sin(angle)
			if err := aperture.DrawApertureSurface(surface, gfxState, centerX + offsetX, centerY + offsetY); err != nil {
				return err
			}
		}
	} else {
		// Else, we can optimize by drawing an arc the thickness of the aperture diameter between the two points, then flashing the
		// aperture at each end to get the endcaps correct
		
		fmt.Printf("Optimized arc\n")
		
		// Draw the stroke, except for the endpoints	
		outerRadius := radius + apertureRadius
		innerRadius := radius - apertureRadius
		arc1StartPointX := centerX + (outerRadius * math.Cos(startAngle))
		arc1StartPointY := centerY + (outerRadius * math.Sin(startAngle))
		arc2StartPointX := centerX + (innerRadius * math.Cos(endAngle))
		arc2StartPointY := centerY + (innerRadius * math.Sin(endAngle))
		surface.MoveTo(arc1StartPointX, arc1StartPointY)
		surface.Arc(centerX, centerY, outerRadius, startAngle, endAngle)
		surface.LineTo(arc2StartPointX, arc2StartPointY)
		surface.ArcNegative(centerX, centerY, innerRadius, endAngle, startAngle)
		surface.LineTo(arc1StartPointX, arc1StartPointY)
		surface.Fill()
		
		// Draw each of the endpoints by flashing the aperture at the endpoints
		startX := centerX + (radius * math.Cos(startAngle))
		startY := centerY + (radius * math.Sin(startAngle))
		endX := centerX + (radius * math.Cos(endAngle))
		endY := centerY + (radius * math.Sin(endAngle))
		// We use the special "no hole" version of the draw, because any holes will
		// have been covered over by the rest of the aperture during the stroke
		aperture.DrawApertureSurfaceNoHole(surface, gfxState, startX, startY)
		aperture.DrawApertureSurfaceNoHole(surface, gfxState, endX, endY)
	}
	
	//TODO: Reset so other draw operations can make their own antialiasing decisions
	surface.SetAntialias(cairo.ANTIALIAS_NONE)
	
	return nil
}

func (aperture *CircleAperture) renderApertureToGraphicsState(gfxState *GraphicsState) {
	// This will render the aperture to a cairo surface the first time it is needed, then
	// cache it in the graphics state.  Subsequent draws of the aperture will used the cached surface
	
	// Construct the surface we're drawing to
	imageSize := int(math.Ceil(aperture.diameter * gfxState.scaleFactor))
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, imageSize, imageSize)
	surface.SetAntialias(cairo.ANTIALIAS_DEFAULT)
	// Scale the surface so we can use unscaled coordinates while rendering the aperture
	surface.Scale(gfxState.scaleFactor, gfxState.scaleFactor)
	// Translate the surface so that the origin is actually the center of the image
	surface.Translate(aperture.diameter / 2.0, aperture.diameter / 2.0)
	
	// Draw the aperture
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}
	
	surface.Arc(0.0, 0.0, aperture.diameter / 2.0, 0, TWO_PI)
	surface.Fill()
	
	// Save the aperture reference before the hole (if any) is rendered, to the no-holes aperture map
	gfxState.renderedAperturesNoHoles[aperture.apertureNumber] = surface
	
	// If present, remove the hole
	if aperture.Hole != nil {
		// If there's a hole, we need to create a copy surface and draw the hole on the copy
		newSurface := copyApertureSurface(surface, gfxState, cairo.ANTIALIAS_DEFAULT, gfxState.scaleFactor, aperture.diameter / 2.0, aperture.diameter / 2.0)
		aperture.DrawHoleSurface(newSurface)
		
		// Then, we save the rendered aperture with the hole to the graphics state
		gfxState.renderedApertures[aperture.apertureNumber] = newSurface
	} else {
		// If there wasn't a hole, we can save the same surface reference as the no-hole aperture in the aperture map
		gfxState.renderedApertures[aperture.apertureNumber] = surface
	}
	
	gfxState.renderedApertures[aperture.apertureNumber].WriteToPNG(fmt.Sprintf("Aperture-%d.png", aperture.apertureNumber))
}

func (aperture *CircleAperture) String() string {
	return fmt.Sprintf("{CA, Diameter: %f, Hole: %v}", aperture.diameter, aperture.Hole)
}