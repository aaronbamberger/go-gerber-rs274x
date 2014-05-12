package gerber_rs274x

type ImageBounds struct {
	xMin float64
	xMax float64
	yMin float64
	yMax float64
	boundsSet bool
}

func newImageBounds() *ImageBounds {
	bounds := new(ImageBounds)
	// For now, everything else is ok as its default (0), but this is here in case that needs to change
	return bounds
}

func (bounds *ImageBounds) updateBounds(xMin float64, xMax float64, yMin float64, yMax float64) {
	// If we haven't seen any bounds yet, we just set the mins and maxes
	// Otherwise, we only update them if they're bigger or smaller
	if !bounds.boundsSet {
		bounds.xMin = xMin
		bounds.xMax = xMax
		bounds.yMin = yMin
		bounds.yMax = yMax
		bounds.boundsSet = true
	} else {
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
}

func (bounds *ImageBounds) updateBoundsAperture(currentX float64, currentY float64, apertureMinSize float64) {
	// Calculate the extents from the current point and the aperture min size
	xMin := currentX - apertureMinSize
	xMax := currentX + apertureMinSize
	yMin := currentY - apertureMinSize
	yMax := currentY + apertureMinSize
	
	// Use update bounds to do the actual work
	bounds.updateBounds(xMin, xMax, yMin, yMax)
}