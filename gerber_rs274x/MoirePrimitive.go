package gerber_rs274x

import (
	"fmt"
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