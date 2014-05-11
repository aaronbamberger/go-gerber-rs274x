package gerber_rs274x

import (
	"fmt"
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

func (primitive *ThermalPrimitive) String() string {
	return fmt.Sprintf("{Thermal, Center (%v %v), Outer Diameter %v, Inner Diameter %v, Gap Thickness %v, Rotation %v}",
						primitive.centerX,
						primitive.centerY,
						primitive.outerDiameter,
						primitive.innerDiameter,
						primitive.gapThickness,
						primitive.rotationAngle)
}