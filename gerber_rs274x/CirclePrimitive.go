package gerber_rs274x

import (
	"fmt"
	cairo "github.com/ungerik/go-cairo"
)

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

func (primitive *CirclePrimitive) GetPrimitiveBounds(env *ExpressionEnvironment) (xMin float64, xMax float64, yMin float64, yMax float64) {
	centerX := primitive.centerX.EvaluateExpression(env)
	centerY := primitive.centerY.EvaluateExpression(env)
	radius := primitive.diameter.EvaluateExpression(env) / 2.0

	return centerX - radius,centerX + radius,centerY - radius,centerY + radius
}

func (primitive *CirclePrimitive) DrawPrimitiveToSurface(surface *cairo.Surface, env *ExpressionEnvironment) error {
	//TODO: Implement
	return nil
}

func (primitive *CirclePrimitive) String() string {
	return fmt.Sprintf("{Circle, Exposure %v, Diameter %v, Center (%v %v)}", primitive.exposure, primitive.diameter, primitive.centerX, primitive.centerY)
}