package gerber_rs274x

import (
	"fmt"
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

func (primitive *CirclePrimitive) String() string {
	return fmt.Sprintf("{Circle, Exposure %v, Diameter %v, Center (%v %v)}", primitive.exposure, primitive.diameter, primitive.centerX, primitive.centerY)
}