package gerber_rs274x

type LowerLeftLinePrimitive struct {
	exposure ApertureMacroExpression
	width ApertureMacroExpression
	height ApertureMacroExpression
	lowerLeftX ApertureMacroExpression
	lowerLeftY ApertureMacroExpression
	rotationAngle ApertureMacroExpression
}

func (primitive *LowerLeftLinePrimitive) AperturePrimitivePlaceholder() {

}

func (primitive *LowerLeftLinePrimitive) ApertureMacroDataBlockPlaceholder() {

}