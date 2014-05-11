package gerber_rs274x

type OutlinePrimitive struct {
	exposure ApertureMacroExpression
	nPoints ApertureMacroExpression
	startX ApertureMacroExpression
	startY ApertureMacroExpression
	subsequentX []ApertureMacroExpression
	subsequentY []ApertureMacroExpression
	rotationAngle ApertureMacroExpression
}

func (primitive *OutlinePrimitive) AperturePrimitivePlaceholder() {

}

func (primitive *OutlinePrimitive) ApertureMacroDataBlockPlaceholder() {

}