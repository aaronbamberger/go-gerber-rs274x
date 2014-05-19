package gerber_rs274x

import (
	cairo "github.com/ungerik/go-cairo"
)

// This controls the number of steps used to render strokes when an optimized draw cannot be used and the aperture
// must be stroked manually (mostly applies to aperture macros and short strokes with standard apertures that have holes,
// aka, strokes less than the shortest radius of the aperture, because then the hole will not be completely obscured)
// The bigger this is, the more accurate non-optimized strokes will be, but the longer they will take.  This has been
// set empirically, it may need to be adjusted later for a different balance of performance and accuracy
const SLOW_DRAWING_STEPS int = 100

type Aperture interface {
	AperturePlaceholder()
	GetApertureNumber() int
	SetHole(hole Hole)
	GetHole() Hole
	GetMinSize(gfxState *GraphicsState) float64
	DrawApertureSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error
	DrawApertureSurfaceNoHole(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error
	DrawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error
	StrokeApertureLinear(surface *cairo.Surface, gfxState *GraphicsState, startX float64, startY float64, endX float64, endY float64) error
	StrokeApertureClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error
	StrokeApertureCounterClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error
	renderApertureToGraphicsState(gfxState *GraphicsState)
}

type Hole interface {
	HolePlaceholder()
	DrawHoleSurface(surface *cairo.Surface) error
}

func renderApertureToSurface(aperture Aperture, surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	return renderApertureToSurfaceHelper(gfxState.renderedApertures, aperture, surface, gfxState, x, y)
}

func renderApertureNoHoleToSurface(aperture Aperture, surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	return renderApertureToSurfaceHelper(gfxState.renderedAperturesNoHoles, aperture, surface, gfxState, x, y)
}

func renderApertureToSurfaceHelper(apertureTable map[int]*cairo.Surface, aperture Aperture, surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	// Draw the aperture
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}
	
	// First, remove the surface scaling (this is because the aperture surfaces are already scaled,
	// and we don't want to scale twice
	surface.Restore()
	
	var renderedAperture *cairo.Surface
	var found bool

	// Try to get the rendered aperture from the graphics state cache.  If it isn't in the cache,
	// we need to actually render it, which will put it in the cache for future use
	if renderedAperture,found = apertureTable[aperture.GetApertureNumber()]; !found {
		// If this is the first use of this aperture, it hasn't been rendered yet,
		// so go ahead and render it before we draw it
		aperture.renderApertureToGraphicsState(gfxState)
		renderedAperture = apertureTable[aperture.GetApertureNumber()]
	}
	
	// Render the aperture to the surface
	surface.MaskSurface(renderedAperture, x * gfxState.scaleFactor, y * gfxState.scaleFactor) // Need to manually scale center coordinates here, because we've removed the surface scaling
	
	// Now, re-apply the surface scaling, so that subsequent draw operations can use the scaling
	surface.Save()
	surface.Scale(gfxState.scaleFactor, gfxState.scaleFactor)
	
	return nil
}

func copyApertureSurface(source *cairo.Surface, gfxState *GraphicsState, antialias cairo.Antialias, scaleFactor float64, xOffset float64, yOffset float64) *cairo.Surface {
	// Create the new surface and initialize settings to be the same as the source surface
	newSurface := cairo.NewSurface(cairo.FORMAT_ARGB32, source.GetWidth(), source.GetHeight())
	newSurface.SetAntialias(antialias)
	
	// Copy the source surface to the new surface
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		newSurface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		newSurface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}
	newSurface.MaskSurface(source, 0.0, 0.0)
	
	// Finally, apply the same transformations as the source surface, so it will behave like the source surface for future drawing operations
	newSurface.Scale(scaleFactor, scaleFactor)
	newSurface.Translate(xOffset, yOffset)
	
	return newSurface
} 
