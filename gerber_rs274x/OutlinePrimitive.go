package gerber_rs274x

type OutlinePrimitive struct {
	exposure Modifier
	nPoints Modifier
	startX Modifier
	startY Modifier
	subsequentX []Modifier
	subsequentY []Modifier
	rotationAngle Modifier
}

func (outline* OutlinePrimitive) PrimitivePlaceholder() {

}