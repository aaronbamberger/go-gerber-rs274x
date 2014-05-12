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
	correctedX := ((x - radius) * gfxState.scaleFactor) + gfxState.xOffset
	correctedY := ((y - radius) * gfxState.scaleFactor) + gfxState.yOffset
	
	// Draw the aperture
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}

	if renderedAperture,found := gfxState.renderedApertures[aperture.apertureNumber]; !found {
		// If this is the first use of this aperture, it hasn't been rendered yet,
		// so go ahead and render it before we draw it
		aperture.renderApertureToGraphicsState(gfxState)
		renderedAperture = gfxState.renderedApertures[aperture.apertureNumber]
		surface.MaskSurface(renderedAperture, correctedX, correctedY)
	} else {
		// Otherwise, just draw the previously rendered aperture
		surface.MaskSurface(renderedAperture, correctedX, correctedY)
	}
	
	return nil
}

func (aperture *CircleAperture) renderApertureToGraphicsState(gfxState *GraphicsState) {
	// This will render the aperture to a cairo surface the first time it is needed, then
	// cache it in the graphics state.  Subsequent draws of the aperture will used the cached surface
	
	scaledDiameter := aperture.diameter * gfxState.scaleFactor
	
	// Construct the surface we're drawing to
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, int(math.Ceil(scaledDiameter)), int(math.Ceil(scaledDiameter)))
	
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