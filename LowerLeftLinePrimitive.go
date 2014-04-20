package gerber_rs274x

type LowerLeftLinePrimitive struct {
	exposure Modifier
	width Modifier
	height Modifier
	lowerLeftX Modifier
	lowerLeftY Modifier
	rotationAngle Modifier
}

func (lowerLeftLine* LowerLeftLinePrimitive) PrimitivePlaceholder() {

}