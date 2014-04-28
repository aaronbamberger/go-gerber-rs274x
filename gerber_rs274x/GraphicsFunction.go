package gerber_rs274x

type Interpolation struct {
	fnCode FunctionCode
	opCode OperationCode
	x float64
	y float64
	i float64
	j float64
	fnCodeValid bool
	opCodeValid bool
	xValid bool
	yValid bool
	iValid bool
	jValid bool
}

type SetCurrentAperture struct {
	apertureNumber int
}

type IgnoreDataBlock struct {
	comment string
}

type GraphicsStateChange struct {
	fnCode FunctionCode
}

func (interpolation* Interpolation) DataBlockPlaceholder() {

}

func (setCurrentAperture* SetCurrentAperture) DataBlockPlaceholder() {

}

func (ignoreDataBlock* IgnoreDataBlock) DataBlockPlaceholder() {

}

func (graphicsStateChange* GraphicsStateChange) DataBlockPlaceholder() {

}