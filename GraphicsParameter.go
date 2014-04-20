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