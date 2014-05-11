package gerber_rs274x

import (
	"fmt"
)

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

func (primitive *CenterLinePrimitive) GetPrimitiveBounds(env *ExpressionEnvironment) (xMin float64, xMax float64, yMin float64, yMax float64) {
	//TODO: Implement
	return 0.0,0.0,0.0,0.0
}

func (primitive *CenterLinePrimitive) String() string {
	return fmt.Sprintf("{Center Line, Exposure %v, Width %v, Height %v, Center (%v %v), Rotation %v}",
						primitive.exposure,
						primitive.width,
						primitive.height,
						primitive.centerX,
						primitive.centerY,
						primitive.rotationAngle)
}