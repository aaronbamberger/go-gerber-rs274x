package gerber_rs274x

type RectangleAperture struct {
	xSize float64
	ySize float64
	Hole
}

func (rectangle* RectangleAperture) AperturePlaceholder() {

}