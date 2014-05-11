package gerber_rs274x

import (
	"fmt"
	"math"
)

type MoirePrimitive struct {
	centerX ApertureMacroExpression
	centerY ApertureMacroExpression
	outerDiameter ApertureMacroExpression
	ringThickness ApertureMacroExpression
	ringGap ApertureMacroExpression
	maxRings ApertureMacroExpression
	crosshairThickness ApertureMacroExpression
	crosshairLength ApertureMacroExpression
	rotationAngle ApertureMacroExpression
}

func (primitive *MoirePrimitive) AperturePrimitivePlaceholder() {

}

func (primitive *MoirePrimitive) ApertureMacroDataBlockPlaceholder() {

}

func (primitive *MoirePrimitive) GetPrimitiveBounds(env *ExpressionEnvironment) (xMin float64, xMax float64, yMin float64, yMax float64) {
	centerX := primitive.centerX.EvaluateExpression(env)
	centerY := primitive.centerY.EvaluateExpression(env)
	ringRadius := primitive.outerDiameter.EvaluateExpression(env) / 2.0
	crosshairRadius := primitive.crosshairLength.EvaluateExpression(env) / 2.0
	maxRadius := math.Max(ringRadius, crosshairRadius) 

	return centerX - maxRadius,centerX + maxRadius,centerY - maxRadius,centerY + maxRadius
}

func (primitive *MoirePrimitive) String() string {
	return fmt.Sprintf("{Moire, Center (%v %v), Outer Diameter %v, Ring Thickness %v, Ring Gap %v, Max Rings %v, Crosshair Thickness %v, CrosshairLength %v, Rotation %v}",
						primitive.centerX,
						primitive.centerY,
						primitive.outerDiameter,
						primitive.ringThickness,
						primitive.ringGap,
						primitive.maxRings,
						primitive.crosshairThickness,
						primitive.crosshairLength,
						primitive.rotationAngle)
}