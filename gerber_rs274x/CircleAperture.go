package gerber_rs274x

type CircleAperture struct {
	diameter float64
	Hole
}

func (circle* CircleAperture) AperturePlaceholder() {

}