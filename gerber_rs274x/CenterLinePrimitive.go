package gerber_rs274x

type CenterLinePrimitive struct {
	exposure ApertureMacroExpression
	width ApertureMacroExpression
	height ApertureMacroExpression
	centerX ApertureMacroExpression
	centerY ApertureMacroExpression
	rotationAngle ApertureMacroExpression
}

func (primitive *CenterLinePrimitive) AperturePrimitivePlaceholder() {

}

func (primitive *CenterLinePrimitive) ApertureMacroDataBlockPlaceholder() {

}