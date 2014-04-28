package gerber_rs274x

type MoirePrimitive struct {
	centerX Modifier
	centerY Modifier
	outerDiameter Modifier
	ringThickness Modifier
	ringGap Modifier
	maxRings Modifier
	crosshairThickness Modifier
	crosshairLength Modifier
	rotationAngle Modifier
}

func (moire* MoirePrimitive) PrimitivePlaceholder() {

}