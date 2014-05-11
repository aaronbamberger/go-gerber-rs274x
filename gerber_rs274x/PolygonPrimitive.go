package gerber_rs274x

type PolygonPrimitive struct {
	exposure ApertureMacroExpression
	nVertices ApertureMacroExpression
	centerX ApertureMacroExpression
	centerY ApertureMacroExpression
	diameter ApertureMacroExpression
	rotationAngle ApertureMacroExpression
}

func (primitive *PolygonPrimitive) AperturePrimitivePlaceholder() {

}

func (primitive *PolygonPrimitive) ApertureMacroDataBlockPlaceholder() {

}