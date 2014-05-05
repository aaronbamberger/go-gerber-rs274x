package gerber_rs274x

import cairo "github.com/ungerik/go-cairo"

type ApertureMacroParameter struct {
	paramCode ParameterCode
	macroName string
	comments []Comment
	primitives []Primitive
	variables map[string]string
}

type Comment struct {
	precedingLine int
	comments []string
}

type Primitive interface {
	PrimitivePlaceholder()
}

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