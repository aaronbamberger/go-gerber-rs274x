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
	correctedX := x - (aperture.xSize / 2.0)
	correctedY := y - (aperture.ySize / 2.0)
	
	return renderApertureToSurface(aperture, surface, gfxState, correctedX, correctedY)
}

func (aperture *RectangleAperture) StrokeApertureLinear(surface *cairo.Surface, gfxState *GraphicsState, startX float64, startY float64, endX float64, endY float64) error {
	radiusX := aperture.xSize / 2.0
	radiusY := aperture.ySize / 2.0
	
	var topLeftX float64
	var topLeftY float64
	var topRightX float64
	var topRightY float64
	var bottomLeftX float64
	var bottomLeftY float64
	var bottomRightX float64
	var bottomRightY float64
	
	if startY < endY {
		topLeftX = startX - radiusX
		topRightX = endX - radiusX
		bottomLeftX = startX + radiusX
		bottomRightX = endX + radiusX
	} else {
		topLeftX = startX + radiusX
		topRightX = endX + radiusX
		bottomLeftX = startX - radiusX
		bottomRightX = endX - radiusX
	}
	
	if startX < endX {
		topLeftY = startY + radiusY
		topRightY = endY + radiusY
		bottomLeftY = startY - radiusY
		bottomRightY = endY - radiusY
	} else {
		topLeftY = endY + radiusY
		topRightY = startY + radiusY
		bottomLeftY = endY - radiusY
		bottomRightY = startY - radiusY
	}

	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}
	
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

	return nil
}

func (aperture *RectangleAperture) StrokeApertureClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *RectangleAperture) StrokeApertureCounterClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *RectangleAperture) renderApertureToGraphicsState(gfxState *GraphicsState) {
	// This will render the aperture to a cairo surface the first time it is needed, then
	// cache it in the graphics state.  Subsequent draws of the aperture will used the cached surface
	
	radiusX := aperture.xSize / 2.0
	radiusY := aperture.ySize / 2.0
	
	// Construct the surface we're drawing to
	imageWidth := int(math.Ceil(aperture.xSize * gfxState.scaleFactor))
	imageHeight := int(math.Ceil(aperture.ySize * gfxState.scaleFactor))
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, imageWidth, imageHeight)
	surface.SetAntialias(cairo.ANTIALIAS_NONE)
	// Scale the surface so we can use unscaled coordinates while rendering the aperture
	surface.Scale(gfxState.scaleFactor, gfxState.scaleFactor)
	// Translate the surface so that the origin is actually the center of the image
	surface.Translate(radiusX, radiusY)
	
	// Draw the aperture
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}
	
	surface.MoveTo(-radiusX, radiusY)
	surface.LineTo(radiusX, radiusY)
	surface.LineTo(radiusX, -radiusY)
	surface.LineTo(-radiusX, -radiusY)
	surface.LineTo(-radiusX, radiusY)
	
	surface.Fill()
	
	// If present, remove the hole
	if aperture.Hole != nil {
		aperture.DrawHoleSurface(surface)
	}
	
	surface.WriteToPNG(fmt.Sprintf("Aperture-%d.png", aperture.apertureNumber))
	
	gfxState.renderedApertures[aperture.apertureNumber] = surface
}

func (aperture *RectangleAperture) String() string {
	return fmt.Sprintf("{RA, X: %f, Y: %f, Hole: %v}", aperture.xSize, aperture.ySize, aperture.Hole)
}