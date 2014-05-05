package gerber_rs274x

import (
	"fmt"
	"math"
	cairo "github.com/ungerik/go-cairo"
)

type ObroundAperture struct {
	apertureNumber int
	xSize float64
	ySize float64
	Hole
}

func (aperture *ObroundAperture) AperturePlaceholder() {

}

func (aperture *ObroundAperture) GetHole() Hole {
	return aperture.Hole
}

func (aperture *ObroundAperture) SetHole(hole Hole) {
	aperture.Hole = hole
}

func (aperture *ObroundAperture) GetMinSize() float64 {
	return math.Min(aperture.xSize / 2.0, aperture.ySize / 2.0)
}

func (aperture *ObroundAperture) DrawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error {
	xRadius := aperture.xSize / 2.0
	yRadius := aperture.ySize / 2.0
	
	xMin := x - xRadius
	xMax := x + xRadius
	yMin := y - yRadius
	yMax := y + yRadius
	
	bounds.updateBounds(xMin, xMax, yMin, yMax)

	return nil
}

func (aperture *ObroundAperture) DrawApertureSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	
	correctedX := ((x - (aperture.xSize / 2.0)) * gfxState.scaleFactor) + gfxState.xOffset
	correctedY := ((y - (aperture.ySize / 2.0)) * gfxState.scaleFactor) + gfxState.yOffset
	
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

func (aperture *ObroundAperture) renderApertureToGraphicsState(gfxState *GraphicsState) {
	// This will render the aperture to a cairo surface the first time it is needed, then
	// cache it in the graphics state.  Subsequent draws of the aperture will used the cached surface
	
	scaledWidth := aperture.xSize * gfxState.scaleFactor
	scaledHeight := aperture.ySize * gfxState.scaleFactor
	
	// Construct the surface we're drawing to
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, int(math.Ceil(scaledWidth)), int(math.Ceil(scaledHeight)))
	
	// Draw the aperture
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}
	
	radiusX := (aperture.xSize / 2.0) * gfxState.scaleFactor
	radiusY := (aperture.ySize / 2.0) * gfxState.scaleFactor
	centerX := scaledWidth / 2.0
	centerY := scaledHeight / 2.0
	
	if aperture.xSize < aperture.ySize {
		rectRadiusY := (scaledHeight - scaledWidth) / 2.0
		surface.MoveTo(centerX - radiusX, centerY - rectRadiusY)
		surface.Arc(centerX, centerY - rectRadiusY, radiusX, math.Pi, 0)
		surface.LineTo(centerX + radiusX, centerY + rectRadiusY)
		surface.Arc(centerX, centerY + rectRadiusY, radiusX, 0, math.Pi)
		surface.LineTo(centerX - radiusX, centerY - rectRadiusY)  
	} else {
		rectRadiusX := (scaledWidth - scaledHeight) / 2.0
		surface.MoveTo(centerX - rectRadiusX, centerY - radiusY)
		surface.LineTo(centerX + rectRadiusX, centerY - radiusY)
		surface.Arc(centerX + rectRadiusX, centerY, radiusY, 1.5 * math.Pi, 0.5 * math.Pi)
		surface.LineTo(centerX - rectRadiusX, centerY + radiusY)
		surface.Arc(centerX - rectRadiusX, centerY, radiusY, 0.5 * math.Pi, 1.5 * math.Pi)
	}
	
	surface.Fill()
	
	// If present, remove the hole
	if aperture.Hole != nil {
		aperture.DrawHoleSurface(surface, gfxState, centerX, centerY)
	}
	
	surface.WriteToPNG(fmt.Sprintf("Aperture-%d.png", aperture.apertureNumber))
	
	gfxState.renderedApertures[aperture.apertureNumber] = surface
}

func (aperture *ObroundAperture) String() string {
	return fmt.Sprintf("{OA, X: %f, Y: %f, Hole: %v}", aperture.xSize, aperture.ySize, aperture.Hole)
}