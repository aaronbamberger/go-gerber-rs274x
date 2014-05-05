package gerber_rs274x

import (
	"math"
)

type ImageBounds struct {
	xMin float64
	xMax float64
	yMin float64
	yMax float64
	smallestApertureSize float64
}

func newImageBounds() *ImageBounds {
	bounds := new(ImageBounds)
	bounds.smallestApertureSize = math.MaxFloat64 // Start this at max double size, so the bounds determination logic works
	// Everything else is ok as its default (0)
	return bounds
}

func (bounds *ImageBounds) updateBounds(xMin float64, xMax float64, yMin float64, yMax float64) {
	if xMin < bounds.xMin {
		bounds.xMin = xMin
	}
	
	if xMax > bounds.xMax {
		bounds.xMax = xMax
	}
	
	if yMin < bounds.yMin {
		bounds.yMin = yMin
	}
	
	if yMax > bounds.yMax {
		bounds.yMax = yMax
	}
}