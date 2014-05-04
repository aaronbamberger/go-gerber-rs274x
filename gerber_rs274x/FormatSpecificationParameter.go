package gerber_rs274x

import (
	"fmt"
	"math"
	"github.com/ajstarks/svgo"
)

type FormatSpecificationParameter struct {
	paramCode ParameterCode
	zeroOmissionMode ZeroOmissionMode
	coordinateNotation CoordinateNotation
	xNumDigits int
	xNumDecimals int
	yNumDigits int
	yNumDecimals int
}

func (formatSpecification *FormatSpecificationParameter) DataBlockPlaceholder() {

}

func (formatSpecification *FormatSpecificationParameter) ProcessDataBlockSVG(svg *svg.SVG, gfxState *GraphicsState) error {
	if gfxState.coordinateNotationSet {
		return fmt.Errorf("Tried to process illegal 2nd format specification parameter")
	}
	
	gfxState.coordinateNotation = formatSpecification.coordinateNotation
	gfxState.coordinateNotationSet = true
	gfxState.drawPrecision = 1.0 / math.Pow10(formatSpecification.xNumDecimals)

	return nil
}

func (fsParam *FormatSpecificationParameter) String() string {
	var zeroOmissionMode string
	var coordinateValueNotation string
	
	switch fsParam.zeroOmissionMode {
		case OMIT_LEADING_ZEROS:
			zeroOmissionMode = "Omit Leading"
			
		case OMIT_TRAILING_ZEROS:
			zeroOmissionMode = "Omit Trailing"
			
		default:
			zeroOmissionMode = "Unknown"
	}
	
	switch fsParam.coordinateNotation {
		case ABSOLUTE_NOTATION:
			coordinateValueNotation = "Absolute"
			
		case INCREMENTAL_NOTATION:
			coordinateValueNotation = "Incremental"
			
		default:
			coordinateValueNotation = "Unknown"
	}
	
	return fmt.Sprintf("{FS, Zero Omission Mode: %s, Coordinate Value Notation: %s, X Int Pos: %d, X Dec Pos: %d, Y Int Pos: %d, Y Dec Pos: %d}",
						zeroOmissionMode,
						coordinateValueNotation,
						fsParam.xNumDigits,
						fsParam.xNumDecimals,
						fsParam.yNumDigits,
						fsParam.yNumDecimals)
}