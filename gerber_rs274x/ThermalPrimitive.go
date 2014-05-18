package gerber_rs274x

import (
	"fmt"
	"math"
	cairo "github.com/ungerik/go-cairo"
)

type ThermalPrimitive struct {
	centerX ApertureMacroExpression
	centerY ApertureMacroExpression
	outerDiameter ApertureMacroExpression
	innerDiameter ApertureMacroExpression
	gapThickness ApertureMacroExpression
	rotationAngle ApertureMacroExpression
}

func (primitive *ThermalPrimitive) AperturePrimitivePlaceholder() {

}

func (primitive *ThermalPrimitive) ApertureMacroDataBlockPlaceholder() {

}

func (primitive *ThermalPrimitive) GetPrimitiveBounds(env *ExpressionEnvironment) (xMin float64, xMax float64, yMin float64, yMax float64) {
	centerX := primitive.centerX.EvaluateExpression(env)
	centerY := primitive.centerY.EvaluateExpression(env)
	radius := primitive.outerDiameter.EvaluateExpression(env) / 2.0

	return centerX - radius,centerX + radius,centerY - radius,centerY + radius
}

func (primitive *ThermalPrimitive) DrawPrimitiveToSurface(surface *cairo.Surface, env *ExpressionEnvironment) error {
	// If there is a rotation angle defined, first check that the center is at the origin
	// (rotations are only allowed if the center is at the origin)
	centerX := primitive.centerX.EvaluateExpression(env)
	centerY := primitive.centerY.EvaluateExpression(env)
	rotation := primitive.rotationAngle.EvaluateExpression(env) * (math.Pi / 180.0)
	
	if rotation != 0.0 && (centerX != 0.0 || centerY != 0.0) {
		return fmt.Errorf("Thermal primitive rotation is only allowed if the center is at the origin")
	}
	
	// Now that we've checked the center, apply the rotation
	surface.Save()
	surface.Rotate(rotation)
	
	surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	
	// Now, draw the thermal
	outerRadius := (primitive.outerDiameter.EvaluateExpression(env) / 2.0)
	innerRadius := (primitive.innerDiameter.EvaluateExpression(env) / 2.0)
	halfGapThickness := (primitive.gapThickness.EvaluateExpression(env) / 2.0)
	
	outerStartX := centerX + halfGapThickness
	outerStartY := centerY + outerRadius
	outerEndX := centerX + outerRadius
	outerEndY := centerY + halfGapThickness
	innerStartX := centerX + innerRadius
	innerStartY := outerEndY
	innerEndX := outerStartX
	innerEndY := centerY + innerRadius
	outerStartAngle := math.Atan2(outerStartY, outerStartX)
	outerEndAngle := math.Atan2(outerEndY, outerEndX)
	innerStartAngle := math.Atan2(innerStartY, innerStartX)
	innerEndAngle := math.Atan2(innerEndY, innerEndX)
	
	// Since the thermal is composed of 4 copies of the same shape, just rotated by 90 degrees,
	// we draw the same shape 4 times, rotating the surface by 90 degrees each time
	
	for i := 0; i < 4; i++ {
		surface.Save()
	
		// Rotate the surface
		surface.Rotate((math.Pi / 2.0) * float64(i))
	
		//Draw one piece of the primitive
		surface.MoveTo(outerStartX, outerStartY)
		surface.ArcNegative(centerX, centerY, outerRadius, outerStartAngle, outerEndAngle)
		surface.LineTo(innerStartX, innerStartY)
		surface.Arc(centerX, centerY, innerRadius, innerStartAngle, innerEndAngle)
		surface.LineTo(outerStartX, outerStartY)
		surface.Fill()
		
		surface.Restore()
	}
	
	// Undo all surface transformations
	surface.Restore()
	
	return nil
}

func (primitive *ThermalPrimitive) String() string {
	return fmt.Sprintf("{Thermal, Center (%v %v), Outer Diameter %v, Inner Diameter %v, Gap Thickness %v, Rotation %v}",
						primitive.centerX,
						primitive.centerY,
						primitive.outerDiameter,
						primitive.innerDiameter,
						primitive.gapThickness,
						primitive.rotationAngle)
}