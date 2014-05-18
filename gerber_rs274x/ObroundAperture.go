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

func (aperture *ObroundAperture) GetApertureNumber() int {
	return aperture.apertureNumber
}

func (aperture *ObroundAperture) GetHole() Hole {
	return aperture.Hole
}

func (aperture *ObroundAperture) SetHole(hole Hole) {
	aperture.Hole = hole
}

func (aperture *ObroundAperture) GetMinSize(gfxState *GraphicsState) float64 {
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
	
	correctedX := (x - (aperture.xSize / 2.0)) * gfxState.scaleFactor
	correctedY := (y - (aperture.ySize / 2.0)) * gfxState.scaleFactor
	
	return renderApertureToSurface(aperture, surface, gfxState, correctedX, correctedY)
}

func (aperture *ObroundAperture) StrokeApertureLinear(surface *cairo.Surface, gfxState *GraphicsState, startX float64, startY float64, endX float64, endY float64) error {
	return nil
}

func (aperture *ObroundAperture) StrokeApertureClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *ObroundAperture) StrokeApertureCounterClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *ObroundAperture) renderApertureToGraphicsState(gfxState *GraphicsState) {
	// This will render the aperture to a cairo surface the first time it is needed, then
	// cache it in the graphics state.  Subsequent draws of the aperture will used the cached surface
	
	radiusX := aperture.xSize / 2.0
	radiusY := aperture.ySize / 2.0
	
	// Construct the surface we're drawing to
	imageWidth := int(math.Ceil(aperture.xSize * gfxState.scaleFactor))
	imageHeight := int(math.Ceil(aperture.ySize * gfxState.scaleFactor))
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, imageWidth, imageHeight)
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
	
	if aperture.xSize < aperture.ySize {
		rectRadiusY := (aperture.ySize - aperture.xSize) / 2.0
		surface.MoveTo(-radiusX, -rectRadiusY)
		surface.Arc(0.0, -rectRadiusY, radiusX, math.Pi, 0)
		surface.LineTo(radiusX, rectRadiusY)
		surface.Arc(0.0, rectRadiusY, radiusX, 0, math.Pi)
		surface.LineTo(-radiusX, -rectRadiusY)  
	} else {
		rectRadiusX := (aperture.xSize - aperture.ySize) / 2.0
		surface.MoveTo(-rectRadiusX, -radiusY)
		surface.LineTo(rectRadiusX, -radiusY)
		surface.Arc(rectRadiusX, 0.0, radiusY, 1.5 * math.Pi, 0.5 * math.Pi)
		surface.LineTo(-rectRadiusX, radiusY)
		surface.Arc(-rectRadiusX, 0.0, radiusY, 0.5 * math.Pi, 1.5 * math.Pi)
	}
	
	surface.Fill()
	
	// If present, remove the hole
	if aperture.Hole != nil {
		aperture.DrawHoleSurface(surface)
	}
	
	surface.WriteToPNG(fmt.Sprintf("Aperture-%d.png", aperture.apertureNumber))
	
	gfxState.renderedApertures[aperture.apertureNumber] = surface
}

func (aperture *ObroundAperture) String() string {
	return fmt.Sprintf("{OA, X: %f, Y: %f, Hole: %v}", aperture.xSize, aperture.ySize, aperture.Hole)
}