package gerber_rs274x

type VectorLinePrimitive struct {
	exposure ApertureMacroExpression
	lineWidth ApertureMacroExpression
	startX ApertureMacroExpression
	startY ApertureMacroExpression
	endX ApertureMacroExpression
	endY ApertureMacroExpression
	rotationAngle ApertureMacroExpression
}

func (primitive *VectorLinePrimitive) AperturePrimitivePlaceholder() {

}

func (primitive *VectorLinePrimitive) ApertureMacroDataBlockPlaceholder() {

}