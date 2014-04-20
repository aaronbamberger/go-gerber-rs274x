package gerber_rs274x

type ObroundAperture struct {
	xSize float64
	ySize float64
	Hole
}

func (obround* ObroundAperture) AperturePlaceholder() {

}