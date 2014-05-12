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

func (aperture *PolygonAperture) renderApertureToGraphicsState(gfxState *GraphicsState) {
	// This will render the aperture to a cairo surface the first time it is needed, then
	// cache it in the graphics state.  Subsequent draws of the aperture will used the cached surface
	
	scaledDiameter := aperture.outerDiameter * gfxState.scaleFactor
	scaledRadius := scaledDiameter / 2.0
	
	// Construct the surface we're drawing to
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, int(math.Ceil(scaledDiameter)), int(math.Ceil(scaledDiameter)))
	
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
	
	// If there's a rotation, convert from the gerber coordinate space
	// to the cairo coordinate space, and apply the rotation
	if (aperture.rotationDegrees != 0.0) {
		// Translate the origin, because we want to rotate about the center of the aperture
		surface.Translate(scaledRadius, scaledRadius)
	
		// Invert the angle (because gerber file and cairo treat signs of angles opposite of each other)
		// and convert to radians
		correctedAngle := aperture.rotationDegrees * (math.Pi / 180.0)
		
		// Perform the rotation
		surface.Rotate(correctedAngle)
		
		// Finally, undo the translation
		surface.Translate(-scaledRadius, -scaledRadius) 
	}
	
	// Move to the first vertex (on the x-axis)
	surface.MoveTo(scaledDiameter, scaledRadius)
	// Draw the edges
	for i := 1; i < aperture.numVertices; i++ {
		xOffset := scaledRadius * math.Cos(float64(i) * vertexAngle)
		yOffset := scaledRadius * math.Sin(float64(i) * vertexAngle)
		surface.LineTo(scaledRadius + xOffset, scaledRadius + yOffset)
	}
	// One final draw to the starting vertex to close the polygon
	surface.LineTo(scaledDiameter, scaledRadius)
	
	surface.Fill()
	
	// Undo any rotations before we draw the hole
	// (holes aren't affected by rotation)
	surface.Restore()
	
	// If present, remove the hole
	if aperture.Hole != nil {
		aperture.DrawHoleSurface(surface, gfxState, scaledRadius, scaledRadius)
	}
	
	surface.WriteToPNG(fmt.Sprintf("Aperture-%d.png", aperture.apertureNumber))
	
	gfxState.renderedApertures[aperture.apertureNumber] = surface
}

func (aperture *PolygonAperture) String() string {
	return fmt.Sprintf("{PA, Diameter: %f, Vertices: %d, Rotation: %f, Hole: %v", aperture.outerDiameter, aperture.numVertices, aperture.rotationDegrees, aperture.Hole)
}