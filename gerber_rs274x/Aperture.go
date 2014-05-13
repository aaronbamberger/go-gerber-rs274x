package gerber_rs274x

import (
	cairo "github.com/ungerik/go-cairo"
)

type Aperture interface {
	AperturePlaceholder()
	GetApertureNumber() int
	SetHole(hole Hole)
	GetHole() Hole
	GetMinSize(gfxState *GraphicsState) float64
	DrawApertureSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error
	DrawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error
	StrokeApertureLinear(surface *cairo.Surface, gfxState *GraphicsState, startX float64, startY float64, endX float64, endY float64) error
	StrokeApertureClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error
	StrokeApertureCounterClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error
	renderApertureToGraphicsState(gfxState *GraphicsState, apertureOffsetX float64, apertureOffsetY float64)
}

type Hole interface {
	HolePlaceholder()
	DrawHoleSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error
}

func renderApertureToSurface(aperture Aperture, surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64, xOffset float64, yOffset float64) error {
	// Draw the aperture
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}

	if renderedAperture,found := gfxState.renderedApertures[aperture.GetApertureNumber()]; !found {
		// If this is the first use of this aperture, it hasn't been rendered yet,
		// so go ahead and render it before we draw it
		aperture.renderApertureToGraphicsState(gfxState, xOffset, yOffset)
		renderedAperture = gfxState.renderedApertures[aperture.GetApertureNumber()]
		surface.MaskSurface(renderedAperture, x, y)
	} else {
		// Otherwise, just draw the previously rendered aperture
		surface.MaskSurface(renderedAperture, x, y)
	}
	
	return nil
}

/*
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
*/