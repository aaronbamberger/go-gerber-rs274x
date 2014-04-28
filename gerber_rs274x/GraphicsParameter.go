package gerber_rs274x

type FormatSpecificationParameter struct {
	paramCode ParameterCode
	zeroOmissionMode ZeroOmissionMode
	coordinateNotation CoordinateNotation
	xNumDigits int
	xNumDecimals int
	yNumDigits int
	yNumDecimals int
}

type LevelPolarityParameter struct {
	paramCode ParameterCode
	polarity Polarity
}

type StepAndRepeatParameter struct {
	paramCode ParameterCode
	xRepeats int
	yRepeats int
	xStepDistance float64
	yStepDistance float64
}

type ModeParameter struct {
	paramCode ParameterCode
	units Units
}

type ApertureDefinitionParameter struct {
	paramCode ParameterCode
	apertureNumber int
	apertureType ApertureType
	aperture Aperture
}

type ImageNameParameter struct {
	paramCode ParameterCode
	name string
}

type ImageRotationParameter struct {
	paramCode ParameterCode
	rotation int
}

type OffsetParameter struct {
	paramCode ParameterCode
	axisAOffset float64
	axisBOffset float64
}

type AxisSelectParameter struct {
	paramCode ParameterCode
	isAXBY bool
}

type ImagePolarityParameter struct {
	paramCode ParameterCode
	polarity Polarity
}

type ScaleFactorParameter struct {
	paramCode ParameterCode
	axisAScale float64
	axisBScale float64
}

type LevelNameParameter struct {
	paramCode ParameterCode
	name string
}

type MirrorImageParameter struct {
	paramCode ParameterCode
	axisAMirror bool
	axisBMirror bool
}

func (apertureMacro* ApertureMacroParameter) DataBlockPlaceholder() {

}

func (formatSpecification* FormatSpecificationParameter) DataBlockPlaceholder() {

}

func (levelPolarity* LevelPolarityParameter) DataBlockPlaceholder() {

}

func (stepAndRepeat* StepAndRepeatParameter) DataBlockPlaceholder() {

}

func (mode* ModeParameter) DataBlockPlaceholder() {

}

func (apertureDefinition* ApertureDefinitionParameter) DataBlockPlaceholder() {

}

func (imageName* ImageNameParameter) DataBlockPlaceholder() {

}

func (imageRotation* ImageRotationParameter) DataBlockPlaceholder() {

}

func (offset* OffsetParameter) DataBlockPlaceholder() {

}

func (axisSelect* AxisSelectParameter) DataBlockPlaceholder() {

}

func (imagePolarity* ImagePolarityParameter) DataBlockPlaceholder() {

}

func (scaleFactor* ScaleFactorParameter) DataBlockPlaceholder() {

}

func (levelName* LevelNameParameter) DataBlockPlaceholder() {

}

func (mirrorImage* MirrorImageParameter) DataBlockPlaceholder() {

}