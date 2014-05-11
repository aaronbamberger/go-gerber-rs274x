package gerber_rs274x

import (
	"fmt"
)

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

func (primitive *LowerLeftLinePrimitive) String() string {
	return fmt.Sprintf("{Lower Left Line, Exposure %v, Width %v, Height %v, Lower Left X %v, Lower Left Y %v, Rotation %v}",
						primitive.exposure,
						primitive.width,
						primitive.height,
						primitive.lowerLeftX,
						primitive.lowerLeftY,
						primitive.rotationAngle)
}