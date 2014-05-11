package gerber_rs274x

type CirclePrimitive struct {
	exposure ApertureMacroExpression
	diameter ApertureMacroExpression
	centerX ApertureMacroExpression
	centerY ApertureMacroExpression
}

func (primitive *CirclePrimitive) AperturePrimitivePlaceholder() {

}

func (primitive *CirclePrimitive) ApertureMacroDataBlockPlaceholder() {

}