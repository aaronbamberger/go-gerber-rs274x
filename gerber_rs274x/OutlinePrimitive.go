package gerber_rs274x

import (
	"fmt"
)

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

func (primitive *OutlinePrimitive) String() string {
	return fmt.Sprintf("{Outline, Exposure %v, Num Points %v, Start X %v, Start Y %v, Subsequent X %v, Subsequent Y %v, Rotation %v}",
						primitive.exposure,
						primitive.nPoints,
						primitive.startX,
						primitive.startY,
						primitive.subsequentX,
						primitive.subsequentY,
						primitive.rotationAngle)
}