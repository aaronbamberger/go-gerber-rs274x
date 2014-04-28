package gerber_rs274x

import (
	"io"
	"bufio"
	"fmt"
	"strconv"
	"regexp"
)

var coordDataBlockRegex *regexp.Regexp
var parameterOrDataBlockRegex *regexp.Regexp
var dataBlockRegex *regexp.Regexp
var dCodeDataBlockRegex *regexp.Regexp
var coordinateDataBlockRegex *regexp.Regexp

type ParseState int

const (
	START_COMMAND ParseState = iota
	PARSING_FUNCTION
	PARSING_COORDINATE_DATA
	PARSING_PARAMETER
	
)

type CoordinateFormat struct {
	numDigits int
	numDecimals int
	suppressTrailingZeros bool
	isSet bool
}

func init() {
	// We compile all regular expressions we'll need for parsing into package global variables, so that we only have to compile
	// them once, not every time they are needed

	coordDataBlockRegex = regexp.MustCompile(`(?:X(?P<xCoord>-?[[:digit:]]*))?(?:Y(?P<yCoord>-?[[:digit:]]*))?(?:I(?P<iOffset>-?[[:digit:]]*))?(?:J(?P<jOffset>-?[[:digit:]]*))?`)
	
	// Regular expression that matches either a parameter block in between "%" characters, or a data block ended by a "*" character
	// NOTE: The parameter matching part is tricky, since parameters can have multiple embedded data blocks (meaning multiple embedded "*" characters)
	// but we still want to match the optional "*" character at the end of the parameter block but not capture it, so we use a non-greedy "*" matching clause
	// inside the capturing expression, but a greedy "*" matching clause at the end of the parameter matching section
	parameterOrDataBlockRegex = regexp.MustCompile("(?:%(?P<paramBlock>(?:[[:alnum:] -\\$&-\\)\\+-/:<-@[-`{-~]*\\**?)*)\\*?%)|(?:(?P<dataBlock>[[:alnum:] -\\$&-\\)\\+-/:<-@[-`{-~]*)\\*)")
	
	dataBlockRegex = regexp.MustCompile(`(?:(?P<fnLetter>G|M)(?P<fnCode>[[:digit:]]{1,2}))?(?P<restOfBlock>[[:alnum:][:punct:] ]*)`)
	
	dCodeDataBlockRegex = regexp.MustCompile(`(?P<restOfBlock>[XYIJ\-[:digit:]]*)(?:D(?P<dCode>[[:digit:]]{1,2}))?`)
	
	coordinateDataBlockRegex = regexp.MustCompile(`(?:X(?P<xCoord>-?[[:digit:]]*))?(?:Y(?P<yCoord>-?[[:digit:]]*))?(?:I(?P<iOffset>-?[[:digit:]]*))?(?:J(?P<jOffset>-?[[:digit:]]*))?`)
}

func ParseGerberFile(in io.Reader) (parsedFile []*Command, err error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	fileString := ""
	for scanner.Scan() {
		fileString += scanner.Text()
	}
	
	if err := scanner.Err(); err != nil {
		return nil,fmt.Errorf("Error encountered while reading file: %v\n", err)
	} 
	
	results := parameterOrDataBlockRegex.FindAllStringSubmatch(fileString, -1)
	
	// Set up the variables we'll need for parsing
	// We'll start with a default size of 100 for now
	// The slice will grow as necessary during parsing
	parsedFile = make([]*Command, 0, 100)
	//var currentCommand Command
	//var coordFormat CoordinateFormat
	
	dataBlocks := make([]DataBlock, 0, 100)
	
	for index,submatch := range results {
		if len(submatch) != 3 {
			return nil,fmt.Errorf("Error (token %d): Parse error on command %v\n", index, submatch)
		}
		
		if len(submatch[1]) > 0 {
			//fmt.Printf("Token %d, Parsed parameter: %s\n", index, submatch[1])
		} else if len(submatch[2]) > 0 {
			if dataBlock,err := parseDataBlock(submatch[2]); err != nil {
				fmt.Printf("Parse Error for block %s: %s\n", submatch[2], err.Error())
			} else {
				dataBlocks = append(dataBlocks, dataBlock)
			}
		} else {
			return nil,fmt.Errorf("Error (token %d): Not parameter or data block: %v\n", index, submatch)
		}
	}
	
	for index,dataBlock := range dataBlocks {
		fmt.Printf("Parsed data block %d: Type: %T, Value: %v\n", index, dataBlock, dataBlock)
	}
	
	return nil,nil
}

func parseDataBlock(dataBlock string) (DataBlock, error) {
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
		return parseNonCommentBlock(parsedDataBlock[0][1], parsedDataBlock[0][2], parsedDataBlock[0][3])
	}
}

func parseNonCommentBlock(fnLetter string, fnCode string, restOfBlock string) (DataBlock, error) {
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
								
								return parseCoordinateDataBlock(parsedDataBlock[0][1], newInterpolation)
							}
						}
					} else {
						// If there was no D code, then this is either a G01, G02, or G03 command by itself with no coordinate data
						// First, make sure there is no coordinate data, because coordinate data without a D code is deprecated
						if len(parsedDataBlock[0][1]) > 0 {
							return nil,fmt.Errorf("Coordinate data without a D code is deprecated and not allowed (Coordinate data: %s)", parsedDataBlock[0][1])
						}
						
						newInterpolation := new(Interpolation)
						newInterpolation.opCodeValid = false
						newInterpolation.xValid = false
						newInterpolation.yValid = false
						newInterpolation.iValid = false
						newInterpolation.jValid = false
						
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

func parseCoordinateDataBlock(restOfBlock string, interpolation *Interpolation) (*Interpolation, error) {
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
		if x,err := strconv.ParseInt(parsedDataBlock[0][1], 10, 32); err != nil {
			return nil,err
		} else {
			//TODO: Scale this appropriately according to the number format
			interpolation.x = float64(x)
			interpolation.xValid = true
		}
	} else {
		interpolation.xValid = false
	}
	
	if len(parsedDataBlock[0][2]) > 0 {
		if y,err := strconv.ParseInt(parsedDataBlock[0][2], 10, 32); err != nil {
			return nil,err
		} else {
			//TODO: Scale this appropriately according to the number format
			interpolation.y = float64(y)
			interpolation.yValid = true
		}
	} else {
		interpolation.yValid = false
	}
	
	if len(parsedDataBlock[0][3]) > 0 {
		if i,err := strconv.ParseInt(parsedDataBlock[0][3], 10, 32); err != nil {
			return nil,err
		} else {
			//TODO: Scale this appropriately according to the number format
			interpolation.i = float64(i)
			interpolation.iValid = true
		}
	} else {
		interpolation.i = 0.0
		interpolation.iValid = true
	}
	
	if len(parsedDataBlock[0][4]) > 0 {
		if j,err := strconv.ParseInt(parsedDataBlock[0][4], 10, 32); err != nil {
			return nil,err
		} else {
			//TODO: Scale this appropriately according to the number format
			interpolation.j = float64(j)
			interpolation.jValid = true
		}
	} else {
		interpolation.j = 0.0
		interpolation.jValid = true
	}
	
	return interpolation,nil
}
