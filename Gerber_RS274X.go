package gerber_rs274x

import (
	"io"
	"bufio"
	"errors"
	"unicode"
	"fmt"
	"strings"
	"strconv"
	"regexp"
)

var coordDataBlockRegex *regexp.Regexp
var parameterOrDataBlockRegex *regexp.Regexp

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

	coordDataBlockRegex = regexp.MustCompile("X?(-?[[:digit:]]*)Y?(-?[[:digit:]]*)I?(-?[[:digit:]]*)J?(-?[[:digit:]]*)")
	
	// Regular expression that matches either a parameter block in between "%" characters, or a data block ended by a "*" character
	// NOTE: The parameter matching part is tricky, since parameters can have multiple embedded data blocks (meaning multiple embedded "*" characters)
	// but we still want to match the optional "*" character at the end of the parameter block but not capture it, so we use a non-greedy "*" matching clause
	// inside the capturing expression, but a greedy "*" matching clause at the end of the parameter matching section
	parameterOrDataBlockRegex = regexp.MustCompile("(?:%((?:[[:alnum:] -\\$&-\\)\\+-/:<-@[-`{-~]*\\**?)*)\\*?%)|(?:([[:alnum:] -\\$&-\\)\\+-/:<-@[-`{-~]*)\\*)")
}

func ParseGerberFile(in io.Reader) (parsedFile []Command, err error) {
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
	parsedFile = make([]Command, 0, 100)
	var currentCommand Command
	var coordFormat CoordinateFormat
	
	for index,submatch := range results {
		if len(submatch) != 3 {
			return nil,fmt.Errorf("Error (token %d): Parse error on command %v\n", index, submatch)
		}
		
		if len(submatch[1]) > 0 {
			fmt.Printf("Token %d, Parsed parameter: %s\n", index, submatch[1])
		} else if len(submatch[2]) > 0 {
			fmt.Printf("Token %d, Parsed data block: %s\n", index, submatch[2])
		} else {
			return nil,fmt.Errorf("Error (token %d): Not parameter or data block: %v\n", index, submatch)
		}
	}

	/*
	reader := (*AsciiReader)(bufio.NewReader(in))
	
	// Set up the variables we'll need for parsing
	// We'll start with a default size of 100 for now
	// The slice will grow as necessary during parsing
	parsedFile = make([]Command, 0, 100)
	var currentCommand Command
	var coordFormat CoordinateFormat
	
	currentState := START_COMMAND
	done := false
	
	for !done {
		if nextChar,err := reader.readAsciiChar(); err != nil {
			if err == io.EOF {
				done = true
				continue
			} else {
				return nil,err
			}
		} else {
			// If we're here, we've successfully read a character from the file
			switch currentState {
			case START_COMMAND:
				switch nextChar {
				case '%':
					currentState = PARSING_PARAMETER
					
				case 'X', 'Y', 'I', 'J', 'D', 'G':
					dataBlock,err := parseFunctionCode(nextChar, &coordFormat, reader)
					
				case 'M':
					
				
				case '\n', '\r':
					continue
				}
			}
		}
	}
	*/
	
	return nil,nil
}

func parseFunctionCode(startChar rune, coordFormat* CoordinateFormat, reader* AsciiReader) (DataBlock, error) {
	interpolation := new(Interpolation)

	switch startChar {
		case 'X':
		
		case 'Y':
		
		case 'I':
		
		case 'J':
		
		case 'D':
		
		case 'G':
			// First, figure out which function code we're talking about
			if firstNum,err := reader.readAsciiChar(); err != nil {
				return nil,err
			} else {
				switch firstNum {
					case '0':
						if secondNum,err := reader.readAsciiChar(); err != nil {
							return nil,err
						} else {
							switch secondNum {
							case '1':
								interpolation.fnCode = LINEAR_INTERPOLATION
								interpolation.fnCodeValid = true
								
							case '2':
								interpolation.fnCode = CIRCULAR_INTERPOLATION_CLOCKWISE
								interpolation.fnCodeValid = true
							
							case '3':
								interpolation.fnCode = CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE
								interpolation.fnCodeValid = false
								
							case '4':
								return parseCommentLine(reader)
								
							default:
								return nil,fmt.Errorf("Unknown function code: G0%c", secondNum)
							}
						}
					
					case '1':
						interpolation.fnCode = LINEAR_INTERPOLATION
						interpolation.fnCodeValid = true
					
					case '2':
						interpolation.fnCode = CIRCULAR_INTERPOLATION_CLOCKWISE
						interpolation.fnCodeValid = true
					
					case '3':
						// For a first character of 3, it's ambiguous whether it's a G03 without the leading 0,
						// or the start of a G36 or G37, so we first have to read the next character to figure it out
						if secondNum,err := reader.readAsciiChar(); err != nil {
							return nil,err
						} else {
							switch secondNum {
								case '6':
									return parseGraphicsStateChange(REGION_MODE_ON, reader)
								
								case '7':
									return parseGraphicsStateChange(REGION_MODE_OFF, reader)
								
								default:
									// If the next character wasn't either a 6 or 7, this was a single character G03,
									// so unread the last character and parse accordingly
									reader.unreadAsciiChar()
									interpolation.fnCode = CIRCULAR_INTERPOLATION_COUNTER_CLOCKWISE
									interpolation.fnCodeValid = true
							}
						}
					
					case '4':
						return parseCommentLine(reader)
					
					case '5':
						if secondNum,err := reader.readAsciiChar(); err != nil {
							return nil,err
						} else {
							switch secondNum {
								case '4':
									return parseGraphicsStateChange(SELECT_APERTURE, reader)
								
								case '5':
									return parseGraphicsStateChange(PREPARE_FOR_FLASH, reader)
								
								default:
									return nil,fmt.Errorf("Unknown function code: G5%c", secondNum)
							}
						}
					
					case '7':
						if secondNum,err := reader.readAsciiChar(); err != nil {
							return nil,err
						} else {
							switch secondNum {
								case '0':
									return parseGraphicsStateChange(SET_UNIT_INCH, reader)
								
								case '1':
									return parseGraphicsStateChange(SET_UNIT_MM, reader)
								
								case '4':
									return parseGraphicsStateChange(SINGLE_QUADRANT_MODE, reader)
								
								case '5':
									return parseGraphicsStateChange(MULTI_QUADRANT_MODE, reader)
								
								default:
									return nil,fmt.Errorf("Unknown function code: G7%c", secondNum)
							}
						}
					
					case '9':
						if secondNum,err := reader.readAsciiChar(); err != nil {
							return nil,err
						} else {
							switch secondNum {
								case '0':
									return parseGraphicsStateChange(SET_NOTATION_ABSOLUTE, reader)
									
								case '1':
									return parseGraphicsStateChange(SET_NOTATION_INCREMENTAL, reader)
								
								default:
									return nil,fmt.Errorf("Unknown function code: G9%c", secondNum)
								
							}
						}
						
					default:
						return nil,fmt.Errorf("Unknown function code prefix: G%c", firstNum)
				
				}
			}
	}
	
	return interpolation,nil
}

func parseCoordinateDataBlock(interpolation* Interpolation, coordFormat* CoordinateFormat, reader* AsciiReader) (*Interpolation, error) {
	if firstChar,err := reader.readAsciiChar(); err != nil {
		return nil,err
	} else {
		switch firstChar {
			case 'X':
				
			
			case 'Y':
			
			case 'I':
			
			case 'J':
			
			case 'D':
				// This is possible if the data block was a function code (G01, G02, or G03) without any corresponding coordinate data
				// In this case, set the coordinate data invalid, parse the operation code, and make sure that the block is closed
				interpolation.xValid = false
				interpolation.yValid = false
				interpolation.iValid = false
				interpolation.jValid = false
				
				return parseOperationCode(interpolation, reader)
				
			default:
				
			
		}
	}
	
	return nil,nil
}

func parseCoordinates(interpolation* Interpolation, coordFormat* CoordinateFormat, reader* AsciiReader) (*Interpolation, error) {
	if firstChar,err := reader.readAsciiChar(); err != nil {
		return nil,err
	} else {
		switch firstChar {
			case 'X':
				
			
			case 'Y':
				interpolation.xValid = false // If we're on Y already, there must be no X
			
			case 'I':
				interpolation.xValid = false // If we're on I already, there must be no X
				interpolation.yValid = false // or Y
			
			case 'J':
				interpolation.xValid = false // If we're on J already, there must be no X
				interpolation.yValid = false // or Y
				interpolation.i = 0          // Also, there must be no I, but since I isn't modal
				interpolation.iValid = true  // we can set it to 0 and mark is as valid
			
			case '*':
			
			default:
		}
	}
	
	return nil,nil
}

func readCoordinate(coordFormat* CoordinateFormat, reader* AsciiReader) (float64, error) {
	// If the coordinate format hasn't already been set, we can't parse any coordinates
	if !coordFormat.isSet {
		return 0.0,errors.New("Can't parse coordinate data before coordinate format (FS Parameter) has been set")
	}
	
	maxDigits := coordFormat.numDigits + coordFormat.numDecimals
	
	// If the coordinate format specifies no digits (for some reason this is allowed by the spec),
	// then check to see if there are any digits provided.  If there are, we'll call it an error
	// (since there are more digits than specified).  Otherwise, we'll unread the character and 
	// return 0
	if maxDigits < 1 {
		if nextChar,err := reader.readAsciiChar(); err != nil {
			return 0.0,err
		} else {
			if (nextChar == '-') || unicode.IsDigit(nextChar) {
				// If it's a negative sign or a digit, it's an error
				return 0.0,fmt.Errorf("Encountered numeric character %c with numeric format specified to have 0 digits", nextChar)
			} else {
				// Unread the character, and return 0
				reader.unreadAsciiChar()
				return 0.0,nil
			}
		}
	}
	
	// Once we've handled the zero size number edge case, we know that we're expecting at least one digit.
	// First, we unconditionally read the 1st character, because it's either a minus sign (-), a digit, or the next
	// coordinate specifier.  If it's the next coordinate specifier, see below for edge case handling.
	// Otherwise, this allows us to only check for digits for the rest of the number.  We then keep accumulating
	// digits until we hit a non-digit character, which we unread, and then finish parsing the number
	// We also keep track of the total number of digits read, and error out if we read more than specified
	// by the coordinate format
	coordAccum := ""
	digitsRead := 0
	if nextChar,err := reader.readAsciiChar(); err != nil {
		return 0.0,err
	} else {
		if !unicode.IsDigit(nextChar) && (nextChar != '-') {
			// This means the format specifier has more than 0 digits of precision,
			// but there were no numeric (or minus sign) characters after the coordinate indicator
			// This is technically possible for a value of 0 (all of the leading or trailing zeros were elided,
			// leaving nothing), and the spec doesn't seem to say anything about this being illegal, so we'll
			// treat it as a valid value of zero.  We need to unread the character just read, to preserve it for
			// the next parse operation
			reader.unreadAsciiChar()
			return 0.0,nil	
		}
		
		// If we just read a digit, increment the digit count (don't increment it for a minus sign)
		if unicode.IsDigit(nextChar) {
			digitsRead++
		}
		coordAccum += string(nextChar)
	}
	
	// Now we can process the rest of the number
	for true {
		if nextChar,err := reader.readAsciiChar(); err != nil {
			return 0.0,err
		} else {
			if !unicode.IsDigit(nextChar) {
				// We're done reading the number, unread the latest character
				// and break out of the loop
				reader.unreadAsciiChar()
				break
			}
			
			if digitsRead == maxDigits {
				// If we've already read the maximum number of digits specified by the
				// format specifier, and we just read another digit, return an error
				return 0.0,fmt.Errorf("Coordinate data block contains number longer than file format specifier allows (%d total digits)", maxDigits)
			}
			
			// If we're here, we can go ahead and append the digit
			coordAccum += string(nextChar)
			digitsRead++
		}
	}
	
	// We now have the string representation of the number in the file.  If the numeric format is set to suppress
	// trailing zeros, we need to correct for this.  If it's set to suppress leading zeros, we don't care, because
	// this won't affect the conversion to floating point
	if coordFormat.suppressTrailingZeros {
		if strings.Contains(coordAccum, "-") && ((len(coordAccum) - 1) < maxDigits) {
			coordAccum += strings.Repeat("0", maxDigits - (len(coordAccum) - 1))	
		} else if !strings.Contains(coordAccum, "-") && (len(coordAccum) < maxDigits) {
			coordAccum += strings.Repeat("0", maxDigits - len(coordAccum))
		}
	}
	
	// Now, we can finally convert to a number and correct for the number of decimal digits
	if number,err := strconv.ParseFloat(coordAccum, 64); err != nil {
		return 0.0,err
	} else {
		if coordFormat.numDecimals == 0 {
			return number,nil
		} else {
			return (number / (10.0 * float64(coordFormat.numDecimals))),nil
		}
	}
}

func parseOperationCode(interpolation* Interpolation, reader* AsciiReader) (*Interpolation, error) {
	// Read until the end-of-block terminator to get the operation code
	if operationCode,err := reader.readAsciiString('*'); err != nil {
		return nil,err
	} else {
		// Slice the end-of-block terminator off the end, then handle
		// the operation code
		switch operationCode[:len(operationCode) - 2] {
			case "01":
				interpolation.opCode = INTERPOLATE_OPERATION
				interpolation.opCodeValid = true
			
			case "02":
				interpolation.opCode = MOVE_OPERATION
				interpolation.opCodeValid = true
			
			case "03":
				interpolation.opCode = FLASH_OPERATION
				interpolation.opCodeValid = true
			
			default:
				return nil,fmt.Errorf("Unknown operation code: %s", operationCode[:len(operationCode) - 2])
		}
		return interpolation,nil
	}
}

func parseCommentLine(reader* AsciiReader) (*IgnoreDataBlock, error) {
	// This is a comment, so read all the way to the block delimiter
	if comment,err := reader.readAsciiString('*'); err != nil {
		return nil,err
	} else {
		// Construct an ignore data block with the comment, slicing off the last character in the string
		// because it contains the end of block character
		return &IgnoreDataBlock{comment[:len(comment)-2]},nil
	}
}

func parseGraphicsStateChange(fnCode FunctionCode, reader* AsciiReader) (*GraphicsStateChange, error) {
	// Make sure the next character is the end-of-block character
	// else it's a parse error
	if nextChar,err := reader.readAsciiChar(); err != nil {
		return nil,err
	} else {
		if nextChar == '*' {
			// We've seen an entire data block, so return the appropriate GraphicsStateChange
			return &GraphicsStateChange{fnCode},nil
		} else {
			var function string
			switch fnCode {
				case REGION_MODE_ON:
					function = "G36"
				
				case REGION_MODE_OFF:
					function = "G37"
				
				case SINGLE_QUADRANT_MODE:
					function = "G74"
				
				case MULTI_QUADRANT_MODE:
					function = "G75"
				
				case SELECT_APERTURE:
					function = "G54"
					
				case PREPARE_FOR_FLASH:
					function = "G55"
				
				case SET_UNIT_INCH:
					function = "G70"
				
				case SET_UNIT_MM:
					function = "G71"
				
				case SET_NOTATION_ABSOLUTE:
					function = "G90"
				
				case SET_NOTATION_INCREMENTAL:
					function = "G91"
				
				case OPTIONAL_STOP:
					function = "M00"
				
				case PROGRAM_STOP:
					function = "M01"
					
				default:
					function = "Unknown"
			}
			return nil,fmt.Errorf("Unexpected character after %s function: %c.  Expected end of block (*)", function, nextChar)
		}
	}	
}

type AsciiReader bufio.Reader

func (asciiReader* AsciiReader) readAsciiChar() (rune, error) {
	// Recover the underlying bufio.Reader*
	reader := (*bufio.Reader)(asciiReader)

	char, size, err := reader.ReadRune()
	
	if err != nil {
		return unicode.ReplacementChar,err
	} else if size != 1 {
		return unicode.ReplacementChar,errors.New("Read non-ASCII character")
	} else {
		return char,nil
	}
}

func (asciiReader* AsciiReader) unreadAsciiChar() error {
	// Recover the underlying bufio.Reader*
	reader := (*bufio.Reader)(asciiReader)
	
	return reader.UnreadRune()
}

func (asciiReader* AsciiReader) readAsciiString(delim byte) (string, error) {
	// Recover the underlying bufio.Reader*
	reader := (*bufio.Reader)(asciiReader)
	
	return reader.ReadString(delim)
}