package gerber_rs274x

import (
	"fmt"
	_ "math"
	"github.com/ungerik/go-cairo"
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
	
	/*
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, 100, 100)
	
	rotation := primitive.rotationAngle.EvaluateExpression(env)
	rotationRadians := rotation * (math.Pi / 180.0)
	surface.Rotate(rotationRadians)
	
	yRadius := primitive.lineWidth.EvaluateExpression(env) / 2.0
	startX := primitive.startX.EvaluateExpression(env)
	startY := primitive.startY.EvaluateExpression(env)
	endX := primitive.endX.EvaluateExpression(env)
	endY := primitive.endY.EvaluateExpression(env)
	lineAngle1 := math.Atan2(startY, startX)
	lineAngle2 := math.Atan2(endY, endX)
	*/
	
	return 0.0,0.0,0.0,0.0
}

func (primitive *VectorLinePrimitive) DrawPrimitiveToSurface(surface *cairo.Surface, env *ExpressionEnvironment) error {
	//TODO: Implement
	return nil
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