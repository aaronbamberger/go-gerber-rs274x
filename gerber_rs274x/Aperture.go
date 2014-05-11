package gerber_rs274x

import (
	cairo "github.com/ungerik/go-cairo"
)

type Aperture interface {
	AperturePlaceholder()
	SetHole(hole Hole)
	GetHole() Hole
	GetMinSize() float64
	DrawApertureSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error
	DrawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error
	
}

type Hole interface {
	HolePlaceholder()
	DrawHoleSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error
}