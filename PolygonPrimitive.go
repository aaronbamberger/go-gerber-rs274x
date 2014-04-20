package gerber_rs274x

type PolygonPrimitive struct {
	exposure Modifier
	nVertices Modifier
	centerX Modifier
	centerY Modifier
	diameter Modifier
	rotationAngle Modifier
}

func (polygon* PolygonPrimitive) PrimitivePlaceholder() {

}