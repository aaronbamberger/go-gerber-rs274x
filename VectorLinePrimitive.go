package gerber_rs274x

type VectorLinePrimitive struct {
	exposure Modifier
	lineWidth Modifier
	startX Modifier
	startY Modifier
	endX Modifier
	endY Modifier
	rotationAngle Modifier
}

func (vectorLine* VectorLinePrimitive) PrimitivePlaceholder() {

}