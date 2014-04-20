package gerber_rs274x

type CirclePrimitive struct {
	exposure Modifier
	diameter Modifier
	centerX Modifier
	centerY Modifier
}

func (circle* CirclePrimitive) PrimitivePlaceholder() {

}