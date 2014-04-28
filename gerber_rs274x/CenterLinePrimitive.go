package gerber_rs274x

type CenterLinePrimitive struct {
	exposure Modifier
	width Modifier
	height Modifier
	centerX Modifier
	centerY Modifier
	rotationAngle Modifier
}

func (centerLine* CenterLinePrimitive) PrimitivePlaceholder() {

}