package gerber_rs274x

import (
	"fmt"
)

type PolygonPrimitive struct {
	exposure ApertureMacroExpression
	nVertices ApertureMacroExpression
	centerX ApertureMacroExpression
	centerY ApertureMacroExpression
	diameter ApertureMacroExpression
	rotationAngle ApertureMacroExpression
}

func (primitive *PolygonPrimitive) AperturePrimitivePlaceholder() {

}

func (primitive *PolygonPrimitive) ApertureMacroDataBlockPlaceholder() {

}

func (primitive *PolygonPrimitive) GetPrimitiveBounds(env *ExpressionEnvironment) (xMin float64, xMax float64, yMin float64, yMax float64) {
	centerX := primitive.centerX.EvaluateExpression(env)
	centerY := primitive.centerY.EvaluateExpression(env)
	radius := primitive.diameter.EvaluateExpression(env) / 2.0

	return centerX - radius,centerX + radius,centerY - radius,centerY + radius
}

func (primitive *PolygonPrimitive) String() string {
	return fmt.Sprintf("{Polygon, Exposure %v, Num Vertices %v, Center (%v %v), Diameter %v, Rotation %v}",
						primitive.exposure,
						primitive.nVertices,
						primitive.centerX,
						primitive.centerY,
						primitive.diameter,
						primitive.rotationAngle)
}