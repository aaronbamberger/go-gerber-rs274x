package gerber_rs274x

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
}

type Hole interface {
	HolePlaceholder()
}