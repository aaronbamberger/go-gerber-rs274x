package gerber-rs274x

type ParameterCode int
type ApertureType int
type FunctionCode int
type OperationCode int
type Polarity int
type ZeroOmissionMode int
type CoordinateNotation int
type Units int

const (
	FS_PARAMETER ParameterCode = iota
	MO_PARAMETER
	AD_PARAMETER
	AM_PARAMETER
	SR_PARAMETER
	LP_PARAMETER
)

const (
	INTERPOLATE_OPERATION OperationCode = iota
	MOVE_OPERATION
	FLASH_OPERATION
)

const (
	LINEAR_INTERPOLATION FunctionCode = iota
	CIRCULAR_INTERPOLATION_CLOCKWISE
	CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE
	IGNORE_DATA_BLOCK
	REGION_MODE_ON
	REGION_MODE_OFF
	SINGLE_QUADRANT_MODE
	MULTI_QUADRANT_MODE
	END_OF_FILE
)

const (
	OMIT_LEADING_ZEROS ZeroOmissionMode = iota
	OMIT_TRAILING_ZEROS
)

const (
	ABSOLUTE_NOTATION CoordinateNotation = iota
	INCREMENTAL_NOTATION
)

const (
	UNITS_IN Units = iota
	UNITS_MM
)

const (
	CIRCLE_APERTURE ApertureType = iota
	RECTANGLE_APERTURE
	OBROUND_APERTURE
	POLYGON_APERTURE
	MACRO_APERTURE
)

const (
	CLEAR_POLARITY Polarity = iota
	DARK_POLARITY
)
	
	


	