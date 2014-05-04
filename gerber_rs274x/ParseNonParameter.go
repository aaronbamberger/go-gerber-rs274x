package gerber_rs274x

import (
	"fmt"
	"strconv"
	"strings"
	"math"
)

func parseDataBlock(dataBlock string, env *ParseEnvironment) (DataBlock, error) {
	parsedDataBlock := dataBlockRegex.FindAllStringSubmatch(dataBlock, -1)
	
	// First, make sure we captured the number of subexpressions we expected
	if len(parsedDataBlock) != 1 {
		return nil,fmt.Errorf("Unable to parse data block %s: error 1", dataBlock)
	} else if len(parsedDataBlock[0]) != 4 {
		return nil,fmt.Errorf("Unable to parse data block %s: error 2", dataBlock)
	}
	
	if parsedDataBlock[0][1] == "G" && (parsedDataBlock[0][2] == "04" || parsedDataBlock[0][2] == "4") {
		// Handle comments as a special case
		return &IgnoreDataBlock{parsedDataBlock[0][3]},nil
	} else {
		// Otherwise, finish processing the data block
		return parseNonCommentBlock(parsedDataBlock[0][1], parsedDataBlock[0][2], parsedDataBlock[0][3], env)
	}
}

func parseNonCommentBlock(fnLetter string, fnCode string, restOfBlock string, env *ParseEnvironment) (DataBlock, error) {
	switch fnLetter {
		case "G", "":
			switch fnCode {
				case "01", "02", "03", "54", "55", "": //NOTE: Codes 54 and 55 are deprecated, the empty function code is for coordinate data blocks with no function
					// Parse the D code out of remainder of the block
					parsedDataBlock := dCodeDataBlockRegex.FindAllStringSubmatch(restOfBlock, -1)
					
					// First, make sure we captured the number of subexpressions we expected
					if len(parsedDataBlock) != 1 {
						return nil,fmt.Errorf("Unable to parse D code from data block %s: error 1", restOfBlock)
					} else if len(parsedDataBlock[0]) != 3 {
						return nil,fmt.Errorf("Unable to parse D code from data block %s: error 2", restOfBlock)
					}
					
					if len(parsedDataBlock[0][2]) > 0 { // This is where the D code was parsed to
						if dCode,err := strconv.ParseInt(parsedDataBlock[0][2], 10, 32); err != nil {
							return nil,err
						} else {
							if dCode >= 10 {
								// If the D code is >= 10, then this is a set aperture command
								return &SetCurrentAperture{int(dCode)},nil
							} else {
								// Else, this is an interpolation, so we set up a new interpolation with the
								// function code and d code, and parse the coordinate data
								newInterpolation := new(Interpolation)
								switch fnCode {
									case "01":
										newInterpolation.fnCode = LINEAR_INTERPOLATION
										newInterpolation.fnCodeValid = true
									
									case "02":
										newInterpolation.fnCode = CIRCULAR_INTERPOLATION_CLOCKWISE
										newInterpolation.fnCodeValid = true
									
									case "03":
										newInterpolation.fnCode = CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE
										newInterpolation.fnCodeValid = true
									
									case "55":
										newInterpolation.fnCode = PREPARE_FOR_FLASH
										newInterpolation.fnCodeValid = true
										
									case "":
										newInterpolation.fnCodeValid = false
									
									//NOTE: We don't have to check for code 54, because it's a deprecated code that optionally precedes
									//an aperture selection D code, which has already been handled above, and can be safely ignored there
								}
								
								switch dCode {
									case 1:
										newInterpolation.opCode = INTERPOLATE_OPERATION
										newInterpolation.opCodeValid = true
										
									case 2:
										newInterpolation.opCode = MOVE_OPERATION
										newInterpolation.opCodeValid = true
									
									case 3:
										newInterpolation.opCode = FLASH_OPERATION
										newInterpolation.opCodeValid = true
									
									default:
										return nil,fmt.Errorf("Unknown D Code: %d", dCode)
								}
								
								return parseCoordinateDataBlock(parsedDataBlock[0][1], newInterpolation, env)
							}
						}
					} else {
						// If there was no D code, then this is either a G01, G02, or G03 command by itself with no coordinate data
						// First, make sure there is no coordinate data, because coordinate data without a D code is deprecated
						if len(parsedDataBlock[0][1]) > 0 {
							return nil,fmt.Errorf("Coordinate data without a D code is deprecated and not allowed (Coordinate data: %s)", parsedDataBlock[0][1])
						}
						
						newInterpolation := new(Interpolation)
						
						switch fnCode {
							case "01":
								newInterpolation.fnCode = LINEAR_INTERPOLATION
								newInterpolation.fnCodeValid = true;
								
							case "02":
								newInterpolation.fnCode = CIRCULAR_INTERPOLATION_CLOCKWISE
								newInterpolation.fnCodeValid = true;
							
							case "03":
								newInterpolation.fnCode = CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE
								newInterpolation.fnCodeValid = true;
							
							default:
								return nil,fmt.Errorf("Illegal function code %s%s by itself", fnLetter, fnCode)
						}
						
						return newInterpolation,nil
					}
				
				case "36":
					return &GraphicsStateChange{REGION_MODE_ON},nil
				
				case "37":
					return &GraphicsStateChange{REGION_MODE_OFF},nil
				
				case "70": //NOTE: Deprecated
					return &GraphicsStateChange{SET_UNIT_INCH},nil
				
				case "71": //NOTE: Deprecated
					return &GraphicsStateChange{SET_UNIT_MM},nil
				
				case "74":
					return &GraphicsStateChange{SINGLE_QUADRANT_MODE},nil
				
				case "75":
					return &GraphicsStateChange{MULTI_QUADRANT_MODE},nil
				
				case "90": //NOTE: Deprecated
					return &GraphicsStateChange{SET_NOTATION_ABSOLUTE},nil
				
				case "91": //NOTE: Deprecated
					return &GraphicsStateChange{SET_NOTATION_INCREMENTAL},nil
				
				default:
					return nil,fmt.Errorf("Error: Unrecognized function code: %s%s", fnLetter, fnCode)
			}
		
		case "M":
			switch fnCode {
				case "00": //NOTE: Deprecated
					return &GraphicsStateChange{PROGRAM_STOP},nil
				
				case "01": //NOTE: Deprecated
					return &GraphicsStateChange{OPTIONAL_STOP},nil
			
				case "02":
					return &GraphicsStateChange{END_OF_FILE},nil
				
				default:
					return nil,fmt.Errorf("Error: Unrecognized function code: %s%s", fnLetter, fnCode)
			}
		
		default:
			return nil,fmt.Errorf("Error: Unrecognized function code: %s%s", fnLetter, fnCode)
	} 
	
	return nil,nil
}

func parseCoordinateDataBlock(restOfBlock string, interpolation *Interpolation, env *ParseEnvironment) (*Interpolation, error) {
	// Make sure the coordinate format has been set
	// It is an error to have a coordinate data block before the coordinate format has been set
	if !env.coordFormat.isSet {
		return nil,fmt.Errorf("Encountered coordinate data block before coordinate format has been set")
	}

	// Parse the rest of the data block
	parsedDataBlock := coordinateDataBlockRegex.FindAllStringSubmatch(restOfBlock, -1)
	
	// First, make sure we captured the number of subexpressions we expected
	if len(parsedDataBlock) != 1 {
		return nil,fmt.Errorf("Unable to parse coordinate data block %s: error 1", restOfBlock)
	} else if len(parsedDataBlock[0]) != 5 {
		return nil,fmt.Errorf("Unable to parse coordinate data block %s: error 2", restOfBlock)
	}
	
	// Parse the coordinate data
	if len(parsedDataBlock[0][1]) > 0 {
		if x,err := parseAndScaleCoordinateData(parsedDataBlock[0][1], env); err != nil {
			return nil,err
		} else {
			interpolation.x = x
			interpolation.xValid = true
		}
	} else {
		interpolation.xValid = false
	}
	
	if len(parsedDataBlock[0][2]) > 0 {
		if y,err := parseAndScaleCoordinateData(parsedDataBlock[0][2], env); err != nil {
			return nil,err
		} else {
			interpolation.y = y
			interpolation.yValid = true
		}
	} else {
		interpolation.yValid = false
	}
	
	if len(parsedDataBlock[0][3]) > 0 {
		if i,err := parseAndScaleCoordinateData(parsedDataBlock[0][3], env); err != nil {
			return nil,err
		} else {
			interpolation.i = i
		}
	} else {
		interpolation.i = 0.0
	}
	
	if len(parsedDataBlock[0][4]) > 0 {
		if j,err := parseAndScaleCoordinateData(parsedDataBlock[0][4], env); err != nil {
			return nil,err
		} else {
			interpolation.j = j
		}
	} else {
		interpolation.j = 0.0
	}
	
	return interpolation,nil
}

func parseAndScaleCoordinateData(coordinateData string, env *ParseEnvironment) (float64, error) {
	// If suppress trailing zeros is enabled, we need to pad out the coordinate string
	// to the max size of the coordinate format with trailing zeros
	if env.coordFormat.suppressTrailingZeros {
		coordMaxLength := env.coordFormat.numDigits + env.coordFormat.numDecimals
		if len(coordinateData) < coordMaxLength {
			coordinateData += strings.Repeat("0", coordMaxLength - len(coordinateData))
		}
	}

	// Next, we parse the string into a double
	if num,err := strconv.ParseFloat(coordinateData, 64); err != nil {
		return 0.0,err
	} else {
		// Now, we scale the number by the coordinate format
		num /= math.Pow10(env.coordFormat.numDecimals)
		return num,nil
	}
}