package gerber_rs274x

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
	IN_PARAMETER // NOTE: Deprecated
	AS_PARAMETER // NOTE: Deprecated
	LN_PARAMETER // NOTE: Deprecated
	IR_PARAMETER // NOTE: Deprecated
	IP_PARAMETER // NOTE: Deprecated
	MI_PARAMETER // NOTE: Deprecated
	OF_PARAMETER // NOTE: Deprecated
	SF_PARAMETER // NOTE: Deprecated
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
	SELECT_APERTURE	// NOTE: Deprecated
	SET_UNIT_INCH // NOTE: Deprecated
	SET_UNIT_MM // NOTE: Deprecated
	SET_NOTATION_INCREMENTAL // NOTE: Deprecated
	SET_NOTATION_ABSOLUTE // NOTE: Deprecated
	OPTIONAL_STOP // NOTE: Deprecated
	PROGRAM_STOP // NOTE: Deprecated
	PREPARE_FOR_FLASH // NOTE: Deprecated
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

type Command struct {
	dataBlocks []DataBlock
}
