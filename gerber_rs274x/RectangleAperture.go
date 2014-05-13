package gerber_rs274x

import (
	"fmt"
	"math"
	cairo "github.com/ungerik/go-cairo"
)

type RectangleAperture struct {
	apertureNumber int
	xSize float64
	ySize float64
	Hole
}

func (aperture *RectangleAperture) AperturePlaceholder() {

}

func (aperture *RectangleAperture) GetApertureNumber() int {
	return aperture.apertureNumber
}

func (aperture *RectangleAperture) GetHole() Hole {
	return aperture.Hole
}

func (aperture *RectangleAperture) SetHole(hole Hole) {
	aperture.Hole = hole
}

func (aperture *RectangleAperture) GetMinSize(gfxState *GraphicsState) float64 {
	return math.Min(aperture.xSize / 2.0, aperture.ySize / 2.0) 
}

func (aperture *RectangleAperture) DrawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error {
	xRadius := aperture.xSize / 2.0
	yRadius := aperture.ySize / 2.0
	
	xMin := x - xRadius
	xMax := x + xRadius
	yMin := y - yRadius
	yMax := y + yRadius
	
	bounds.updateBounds(xMin, xMax, yMin, yMax)
	
	return nil
}

func (aperture *RectangleAperture) DrawApertureSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	correctedX := ((x - (aperture.xSize / 2.0)) * gfxState.scaleFactor) + gfxState.xOffset
	correctedY := ((y - (aperture.ySize / 2.0)) * gfxState.scaleFactor) + gfxState.yOffset
	
	return renderApertureToSurface(aperture, surface, gfxState, correctedX, correctedY, 0.0, 0.0)
}

func (aperture *RectangleAperture) StrokeApertureLinear(surface *cairo.Surface, gfxState *GraphicsState, startX float64, startY float64, endX float64, endY float64) error {
	return nil
}

func (aperture *RectangleAperture) StrokeApertureClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *RectangleAperture) StrokeApertureCounterClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *RectangleAperture) renderApertureToGraphicsState(gfxState *GraphicsState, apertureOffsetX float64, apertureOffsetY float64) {
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
	
	
	surface.MoveTo(0, 0)
	surface.LineTo(scaledWidth, 0)
	surface.LineTo(scaledWidth, scaledHeight)
	surface.LineTo(0, scaledHeight)
	surface.LineTo(0, 0)
	
	surface.Fill()
	
	// If present, remove the hole
	if aperture.Hole != nil {
		centerX := scaledWidth / 2.0
		centerY := scaledHeight / 2.0
		aperture.DrawHoleSurface(surface, gfxState, centerX, centerY)
	}
	
	surface.WriteToPNG(fmt.Sprintf("Aperture-%d.png", aperture.apertureNumber))
	
	gfxState.renderedApertures[aperture.apertureNumber] = surface
}

func (aperture *RectangleAperture) String() string {
	return fmt.Sprintf("{RA, X: %f, Y: %f, Hole: %v}", aperture.xSize, aperture.ySize, aperture.Hole)
}