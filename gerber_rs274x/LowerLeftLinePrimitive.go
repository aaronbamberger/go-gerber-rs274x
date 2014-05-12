package gerber_rs274x

import (
	"fmt"
	cairo "github.com/ungerik/go-cairo"
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

func (primitive *LowerLeftLinePrimitive) GetPrimitiveBounds(env *ExpressionEnvironment) (xMin float64, xMax float64, yMin float64, yMax float64) {
	//TODO: Implement
	return 0.0,0.0,0.0,0.0
}

func (primitive *LowerLeftLinePrimitive) DrawPrimitiveToSurface(surface *cairo.Surface, env *ExpressionEnvironment, scaleFactor float64, xOffset float64, yOffset float64) error {
	//TODO: Implement
	return nil
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