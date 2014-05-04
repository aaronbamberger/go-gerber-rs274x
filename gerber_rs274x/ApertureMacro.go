package gerber_rs274x

import "github.com/ajstarks/svgo"

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
	DrawApertureSVG(svg *svg.SVG, gfxState *GraphicsState, x float64, y float64) error
}

type Hole interface {
	HolePlaceholder()
}