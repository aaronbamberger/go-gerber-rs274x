package gerber_rs274x

import (
	"fmt"
	"math"
	cairo "github.com/ungerik/go-cairo"
)

type PolygonAperture struct {
	apertureNumber int
	outerDiameter float64
	numVertices int
	rotationDegrees float64
	Hole
}

func (aperture *PolygonAperture) AperturePlaceholder() {
	
}

func (aperture *PolygonAperture) GetApertureNumber() int {
	return aperture.apertureNumber
}

func (aperture *PolygonAperture) GetHole() Hole {
	return aperture.Hole
}

func (aperture *PolygonAperture) SetHole(hole Hole) {
	aperture.Hole = hole
}

func (aperture *PolygonAperture) GetMinSize(gfxState *GraphicsState) float64 {
	return aperture.outerDiameter / 2.0
}

func (aperture *PolygonAperture) DrawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error {
	radius := aperture.outerDiameter / 2.0
	
	xMin := x - radius
	xMax := x + radius
	yMin := y - radius
	yMax := y + radius
	
	bounds.updateBounds(xMin, xMax, yMin, yMax)

	return nil
}

func (aperture *PolygonAperture) DrawApertureSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {

	radius := aperture.outerDiameter / 2.0
	correctedX := x - radius
	correctedY := y - radius
	
	return renderApertureToSurface(aperture, surface, gfxState, correctedX, correctedY) 
}

func (aperture *PolygonAperture) StrokeApertureLinear(surface *cairo.Surface, gfxState *GraphicsState, startX float64, startY float64, endX float64, endY float64) error {
	return nil
}

func (aperture *PolygonAperture) StrokeApertureClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *PolygonAperture) StrokeApertureCounterClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *PolygonAperture) renderApertureToGraphicsState(gfxState *GraphicsState) {
	// This will render the aperture to a cairo surface the first time it is needed, then
	// cache it in the graphics state.  Subsequent draws of the aperture will used the cached surface
	radius := aperture.outerDiameter / 2.0
	
	// Construct the surface we're drawing to
	imageSize := int(math.Ceil(aperture.outerDiameter * gfxState.scaleFactor))
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, imageSize, imageSize)
	// Scale the surface so we can use unscaled coordinates while rendering the aperture
	surface.Scale(gfxState.scaleFactor, gfxState.scaleFactor)
	// Translate the surface so that the origin is actually the center of the image
	surface.Translate(radius, radius)
	
	// Draw the aperture
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}
	
	vertexAngle := (2.0 * math.Pi) / float64(aperture.numVertices)
	
	// Save the current surface state so we can undo any
	// rotations that we apply
	surface.Save()
	
	// If there's a rotation, then apply it
	if (aperture.rotationDegrees != 0.0) {
		// Convert the angle to radians
		correctedAngle := aperture.rotationDegrees * (math.Pi / 180.0)
		
		// Perform the rotation
		surface.Rotate(correctedAngle) 
	}
	
	// Move to the first vertex (on the x-axis)
	surface.MoveTo(radius, 0.0)
	// Draw the edges
	for i := 1; i < aperture.numVertices; i++ {
		xOffset := radius * math.Cos(float64(i) * vertexAngle)
		yOffset := radius * math.Sin(float64(i) * vertexAngle)
		surface.LineTo(xOffset, yOffset)
	}
	// One final draw to the starting vertex to close the polygon
	surface.LineTo(radius, 0.0)
	
	surface.Fill()
	
	// Undo any rotations before we draw the hole
	// (holes aren't affected by rotation)
	surface.Restore()
	
	// If present, remove the hole
	if aperture.Hole != nil {
		aperture.DrawHoleSurface(surface)
	}
	
	surface.WriteToPNG(fmt.Sprintf("Aperture-%d.png", aperture.apertureNumber))
	
	gfxState.renderedApertures[aperture.apertureNumber] = surface
}

func (aperture *PolygonAperture) String() string {
	return fmt.Sprintf("{PA, Diameter: %f, Vertices: %d, Rotation: %f, Hole: %v", aperture.outerDiameter, aperture.numVertices, aperture.rotationDegrees, aperture.Hole)
}