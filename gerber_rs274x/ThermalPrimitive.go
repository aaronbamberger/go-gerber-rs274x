package gerber_rs274x

type ThermalPrimitive struct {
	centerX Modifier
	centerY Modifier
	outerDiameter Modifier
	innerDiameter Modifier
	gapThickness Modifier
	rotationAngle Modifier
}

func (thermal* ThermalPrimitive) PrimitivePlaceholder() {

}