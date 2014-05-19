package gerber_rs274x

import (
	"math"
	cairo "github.com/ungerik/go-cairo"
)

type GraphicsState struct {
	currentAperture int
	currentQuadrantMode FunctionCode
	currentInterpolationMode FunctionCode
	currentX float64
	currentY float64
	currentLevelPolarity Polarity
	regionModeOn bool
	xImageSize int
	yImageSize int
	fileComplete bool
	coordinateNotation CoordinateNotation
	filePrecision float64
	ScalingParms
	
	// As we encounter aperture definitions, we save them
	// for later use while drawing
	apertures map[int]Aperture
	// We also need to remember aperture macro definitions, so that we can recall them when they are
	// referenced in aperture definition parameters
	apertureMacros map[string][]ApertureMacroDataBlock
	// The first time an aperture is rendered, we render it to a cairo surface
	// Then, we can just look up the rendered aperture the next time we need it
	// This should provide for some optimization, since the same aperture will
	// get used over and over to stroke a path
	renderedApertures map[int]*cairo.Surface
	// We also need to save apertures rendered without their holes (if they have holes), for use in certain
	// optimized stroke drawing routines.  If a particular aperture doesn't have a hole, a reference to the same
	// surface as is stored in the renderedApertures map is stored here, to save on memory.  Else, a new rendered
	// surface without the hole is stored here
	renderedAperturesNoHoles map[int]*cairo.Surface
	
	// Some of these default to undefined,
	// so we also need to keep track of when they get defined
	apertureSet bool
	quadrantModeSet bool
	interpolationModeSet bool
	coordinateNotationSet bool
}

func newGraphicsState(bounds *ImageBounds, xImageSize int, yImageSize int) *GraphicsState {
	graphicsState := new(GraphicsState)
	
	graphicsState.currentLevelPolarity = DARK_POLARITY
	graphicsState.apertures = make(map[int]Aperture, 10) // Start with an initial capacity of 10 apertures, will grow as needed
	graphicsState.renderedApertures = make(map[int]*cairo.Surface, 10) // Same as above
	graphicsState.renderedAperturesNoHoles = make(map[int]*cairo.Surface, 10) // Same as above
	graphicsState.apertureMacros = make(map[string][]ApertureMacroDataBlock, 10) // Same as above
	
	if bounds != nil {
		// If bounds are provided, compute the necessary scaling information
		xSpan := bounds.xMax - bounds.xMin
		ySpan := bounds.yMax - bounds.yMin
		
		// Build 5% margin on each side into the scaling
		xMargin := float64(xImageSize) * 0.1
		yMargin := float64(yImageSize) * 0.1
		
		// Compute the appropriate scaling factor
		xScale := (float64(xImageSize) - xMargin) / xSpan
		yScale := (float64(yImageSize) - yMargin) / ySpan
		graphicsState.scaleFactor = math.Min(xScale, yScale)
		
		// Compute offsets to apply to all coordinates to start them at zero and account for margins
		graphicsState.xOffset = -(bounds.xMin * graphicsState.scaleFactor) + (xMargin / 2.0)
		graphicsState.yOffset = -(bounds.yMin * graphicsState.scaleFactor) + (yMargin / 2.0)
	}
	
	// All other settings are fine with their go defaults
	// Current aperture: Doesn't matter since it's undefined by default
	// Current quadrant mode: Doesn't matter since it's undefined by default
	// Current interpolation mode: Doesn't matter since it's undefined by default
	// Coordinate notation: Doesn't matter since it's undefined by default
	// Current x: 0 is correct
	// Current y: 0 is correct
	// Region mode on: false is correct
	// Aperture set: false is correct
	// Quadrant mode set: false is correct
	// Interpolation mode set: false is correct
	// Region mode on: false is correct
	// File complete: false is correct
	// Coordinate notation set: false is correct
	
	return graphicsState 
}

func (gfxState *GraphicsState) updateCurrentCoordinate(newX float64, newY float64) {
	gfxState.currentX = newX
	gfxState.currentY = newY
}

func (gfxState *GraphicsState) releaseRenderedSurfaces() {

	// NOTE: For the purposes of this iteration, we assume that the keyset of renderedApertures is a strict superset
	// of the keyset of renderedAperturesNoHoles (which should always be the case).  If this assumption is broken,
	// the iteration may miss releasing some surfaces in the renderedAperturesNoHoles map
	for apertureNumber,surface1 := range gfxState.renderedApertures {
		if surface2,found := gfxState.renderedAperturesNoHoles[apertureNumber]; found {
			// This aperture has been rendered into both the holes and no-holes map
			// Next, we need to check whether both surfaces actually refer to the same object,
			// and if so, only release it once
			if surface1 == surface2 {
				surface1.Finish()
				surface1.Destroy()
			} else {
				// Else, they're two separate objects and we must release them both
				surface1.Finish()
				surface1.Destroy()
				surface2.Finish()
				surface2.Destroy()
			}
		} else {
			// This aperture was only rendered into the holes map
			surface1.Finish()
			surface1.Destroy()
		}
	}
}