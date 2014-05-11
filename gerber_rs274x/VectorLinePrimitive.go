package gerber_rs274x

import (
	"fmt"
)

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

func (primitive *VectorLinePrimitive) String() string {
	return fmt.Sprintf("{Vector Line, Exposure %v, Line Width %v, Start (%v %v), End (%v %v), Rotation %v}",
						primitive.exposure,
						primitive.lineWidth,
						primitive.startX,
						primitive.startY,
						primitive.endX,
						primitive.endY,
						primitive.rotationAngle)
}