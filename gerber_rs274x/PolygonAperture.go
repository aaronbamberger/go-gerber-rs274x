package gerber_rs274x

type PolygonAperture struct {
	outerDiameter float64
	numVertices int
	rotationDegrees float64
	Hole
}

func (polygon* PolygonAperture) PolygonPlaceholder() {

}