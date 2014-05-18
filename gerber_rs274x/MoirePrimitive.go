package gerber_rs274x

import (
	"fmt"
	"math"
	cairo "github.com/ungerik/go-cairo"
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

func (primitive *MoirePrimitive) DrawPrimitiveToSurface(surface *cairo.Surface, env *ExpressionEnvironment) error {
	// If there is a rotation angle defined, first check that the center is at the origin
	// (rotations are only allowed if the center is at the origin)
	centerX := primitive.centerX.EvaluateExpression(env)
	centerY := primitive.centerY.EvaluateExpression(env)
	rotation := primitive.rotationAngle.EvaluateExpression(env) * (math.Pi / 180.0)
	
	if rotation != 0.0 && (centerX != 0.0 || centerY != 0.0) {
		return fmt.Errorf("Moire primitive rotation is only allowed if the center is at the origin")
	}
	
	// Now that we've checked the center, first apply a translation to account for the offset,
	// then apply the rotation
	surface.Save()
	surface.Rotate(rotation)
	
	surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	
	// Start drawing the rings
	maxRings := int(primitive.maxRings.EvaluateExpression(env))
	radius := (primitive.outerDiameter.EvaluateExpression(env) / 2.0)
	thickness := primitive.ringThickness.EvaluateExpression(env)
	gap := primitive.ringGap.EvaluateExpression(env)
	for ring := 0; ring < maxRings; ring++ {
		outerRadius := radius - ((thickness + gap) * float64(ring))
		innerRadius := outerRadius - thickness
		
		// Draw the outer portion of the ring
		surface.Arc(centerX, centerY, outerRadius, 0.0, 2.0 * math.Pi)
		
		if innerRadius > 0.0 {
			// Draw the inner portion of the ring
			surface.Arc(centerX, centerY, innerRadius, 0.0, 2.0 * math.Pi)
			surface.Fill()
		} else {
			// We've reached the center, so fill the surface and break out of the loop
			surface.Fill()
			break
		}
	}
	
	// Now, draw the crosshair
	crosshairHalfLength := (primitive.crosshairLength.EvaluateExpression(env) / 2.0)
	crosshairHalfThickness := (primitive.crosshairThickness.EvaluateExpression(env) / 2.0)
	horzLeftX := centerX - crosshairHalfLength
	horzRightX := centerX + crosshairHalfLength
	horzTopY := centerY + crosshairHalfThickness
	horzBottomY := centerY - crosshairHalfThickness
	vertLeftX := centerX - crosshairHalfThickness
	vertRightX := centerX + crosshairHalfThickness
	vertTopY := centerY + crosshairHalfLength
	vertBottomY := centerY - crosshairHalfLength
	// Horizontal crosshair portion
	surface.MoveTo(horzLeftX, horzTopY)
	surface.LineTo(horzRightX, horzTopY)
	surface.LineTo(horzRightX, horzBottomY)
	surface.LineTo(horzLeftX, horzBottomY)
	surface.LineTo(horzLeftX, horzTopY)
	surface.Fill()
	// Vertical crosshair portion
	surface.MoveTo(vertLeftX, vertTopY)
	surface.LineTo(vertRightX, vertTopY)
	surface.LineTo(vertRightX, vertBottomY)
	surface.LineTo(vertLeftX, vertBottomY)
	surface.LineTo(vertLeftX, vertTopY)
	surface.Fill()
	
	fmt.Printf("Horizontal Top Left (%f %f) Bottom Right (%f %f)\n", horzLeftX, horzTopY, horzRightX, horzBottomY)
	fmt.Printf("Vertical Top Left (%f %f) Bottom Right (%f %f)\n", vertLeftX, vertTopY, vertRightX, vertBottomY)
	
	// Finally, undo the transformations to the surface
	surface.Restore()
	
	return nil
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