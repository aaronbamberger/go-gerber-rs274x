package gerber_rs274x

import (
	"fmt"
	cairo "github.com/ungerik/go-cairo"
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

func (primitive *OutlinePrimitive) GetPrimitiveBounds(env *ExpressionEnvironment) (xMin float64, xMax float64, yMin float64, yMax float64) {
	//TODO: Implement
	return 0.0,0.0,0.0,0.0
}

func (primitive *OutlinePrimitive) DrawPrimitiveToSurface(surface *cairo.Surface, env *ExpressionEnvironment) error {
	//TODO: Implement
	return nil
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