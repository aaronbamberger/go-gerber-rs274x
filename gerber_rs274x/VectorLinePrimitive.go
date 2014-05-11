package gerber_rs274x

import (
	"fmt"
	_ "github.com/ungerik/go-cairo"
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

func (primitive *VectorLinePrimitive) GetPrimitiveBounds(env *ExpressionEnvironment) (xMin float64, xMax float64, yMin float64, yMax float64) {
	//TODO: Implement
	
	//surface := cairo.NewSurface(cairo.FORMAT_ARGB32, 100, 100)
	
	return 0.0,0.0,0.0,0.0
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