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

func (primitive *PolygonPrimitive) String() string {
	return fmt.Sprintf("{Polygon, Exposure %v, Num Vertices %v, Center (%v %v), Diameter %v, Rotation %v}",
						primitive.exposure,
						primitive.nVertices,
						primitive.centerX,
						primitive.centerY,
						primitive.diameter,
						primitive.rotationAngle)
}