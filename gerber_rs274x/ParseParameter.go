package gerber_rs274x

import (
	"fmt"
	"strconv"
	"strings"
)

func parseParameter(parameter string, env *ParseEnvironment) (DataBlock, error) {
	// All parameter blocks must have at least 3 characters (the two character parameter code, and at least one character of arguments)
	// So we check for at least that length here, so we can slice to at least the third character below
	if len(parameter) < 3 {
		return nil,fmt.Errorf("Error: Unrecognized parameter string %s", parameter)
	}

	switch parameter[0:2] {
		case "FS":
			newFSParam := new(FormatSpecificationParameter)
			newFSParam.paramCode = FS_PARAMETER
			return parseFSParameter(newFSParam, parameter[2:], env)
		
		case "MO":
			newMOParam := new(ModeParameter)
			newMOParam.paramCode = MO_PARAMETER
			return parseMOParameter(newMOParam, parameter[2:], env)		
		
		case "AD":
			newADParam := new(ApertureDefinitionParameter)
			newADParam.paramCode = AD_PARAMETER
			return parseADParameter(newADParam, parameter[2:], env)
		
		case "AM":
			newAMParam := new(ApertureMacroParameter)
			newAMParam.paramCode = AM_PARAMETER
			return parseAMParameter(newAMParam, parameter[2:])
		
		case "SR":
			newSRParam := new(StepAndRepeatParameter)
			newSRParam.paramCode = SR_PARAMETER
			return parseSRParameter(newSRParam, parameter[2:])
		
		case "LP":
			newLPParam := new(LevelPolarityParameter)
			newLPParam.paramCode = LP_PARAMETER
			return parseLPParameter(newLPParam, parameter[2:])
			
		
		default:
			return nil,fmt.Errorf("Error: Unrecognized parameter code %s", parameter[0:2])
	}
	
	return nil,nil
}

func parseFSParameter(fsParameter *FormatSpecificationParameter, restOfParameter string, env *ParseEnvironment) (DataBlock, error) {
	parsedFS := fsParameterRegex.FindAllStringSubmatch(restOfParameter, -1)
	
	// Make sure we haven't already seen an FS parameter
	// It's only legal to have one FS parameter per file
	if env.coordFormat.isSet {
		return nil,fmt.Errorf("Illegal 2nd FS parameter encountered")
	}
	
	// Make sure we captured the number of subexpressions we expected
	if len(parsedFS) != 1 {
		return nil,fmt.Errorf("Unable to parse FS Parameter %s: error 1", restOfParameter)
	} else if len(parsedFS[0]) != 7 {
		return nil,fmt.Errorf("Unable to parse FS Parameter %s: error 2", restOfParameter)
	}
	
	// Parse zero omission mode
	switch parsedFS[0][1] {
		case "L":
			fsParameter.zeroOmissionMode = OMIT_LEADING_ZEROS
		
		case "T":
			fsParameter.zeroOmissionMode = OMIT_TRAILING_ZEROS
		
		default:
			return nil,fmt.Errorf("Unknown zero omission mode %s", parsedFS[0][1])
	}
	
	// Parse coordinate value notation
	switch parsedFS[0][2] {
		case "A":
			fsParameter.coordinateNotation = ABSOLUTE_NOTATION
		
		case "I":
			fsParameter.coordinateNotation = INCREMENTAL_NOTATION
		
		default:
			return nil,fmt.Errorf("Unknown coordinate value notation %s", parsedFS[0][2])
	}
	
	// Parse X number of integer positions
	if len(parsedFS[0][3]) > 0 {
		if xIntPos,err := strconv.ParseInt(parsedFS[0][3], 10, 32); err != nil {
			return nil,err
		} else {
			if (xIntPos < 0) || (xIntPos > 7) {
				return nil,fmt.Errorf("X coordinate number of integer positions must be between 0 and 7.  Received %d", xIntPos)
			} else {
				fsParameter.xNumDigits = int(xIntPos)
			}
		}
	} else {
		return nil,fmt.Errorf("Missing X coordinate number of integer positions")
	}
	
	// Parse X number of decimal positions
	if len(parsedFS[0][4]) > 0 {
		if xDecPos,err := strconv.ParseInt(parsedFS[0][4], 10, 32); err != nil {
			return nil,err
		} else {
			if (xDecPos < 0) || (xDecPos > 7) {
				return nil,fmt.Errorf("X coordinate number of decimal positions must be between 0 and 7.  Received %d", xDecPos)
			} else {
				fsParameter.xNumDecimals = int(xDecPos)
			}
		}
	} else {
		return nil,fmt.Errorf("Missing X coordinate number of decimal positions")
	}
	
	// Parse Y number of integer positions
	if len(parsedFS[0][5]) > 0 {
		if yIntPos,err := strconv.ParseInt(parsedFS[0][5], 10, 32); err != nil {
			return nil,err
		} else {
			if (yIntPos < 0) || (yIntPos > 7) {
				return nil,fmt.Errorf("Y coordinate number of integer positions must be between 0 and 7.  Received %d", yIntPos)
			} else {
				fsParameter.yNumDigits = int(yIntPos)
			}
		}
	} else {
		return nil,fmt.Errorf("Missing Y coordinate number of integer positions")
	}
	
	// Parse Y number of decimal positions
	if len(parsedFS[0][6]) > 0 {
		if yDecPos,err := strconv.ParseInt(parsedFS[0][6], 10, 32); err != nil {
			return nil,err
		} else {
			if (yDecPos < 0) || (yDecPos > 7) {
				return nil,fmt.Errorf("Y coordinate number of decimal positions must be between 0 and 7.  Received %d", yDecPos)
			} else {
				fsParameter.yNumDecimals = int(yDecPos)
			}
		}
	} else {
		return nil,fmt.Errorf("Missing Y coordinate number of decimal positions")
	}
	
	// Per the spec, the X and Y coordinate formats need to be the same, so enforce that here
	if (fsParameter.xNumDigits != fsParameter.yNumDigits) || (fsParameter.xNumDecimals != fsParameter.yNumDecimals) {
		return nil,fmt.Errorf("X and Y coordinate formats must match.  Received X=(%d %d) Y=(%d %d)", fsParameter.xNumDigits, fsParameter.xNumDecimals, fsParameter.yNumDigits, fsParameter.yNumDecimals)
	}
	
	// If we're here, we've succesfully parsed the FS parameter.  Update the parse environment
	env.coordFormat.numDigits = fsParameter.xNumDigits // We've already verified that x and y have the same coordinate format, so we can just use x here
	env.coordFormat.numDecimals = fsParameter.xNumDecimals // Same as above
	env.coordFormat.suppressTrailingZeros = (fsParameter.zeroOmissionMode == OMIT_TRAILING_ZEROS) 
	env.coordFormat.isSet = true

	return fsParameter,nil
}

func parseMOParameter(moParameter *ModeParameter, restOfParameter string, env *ParseEnvironment) (DataBlock, error) {
	// Make sure we haven't already seen an MO parameter
	// It's only legal to have one MO parameter per file
	if env.unitsSet {
		return nil,fmt.Errorf("Illegal 2nd MO parameter encountered")
	}

	if len(restOfParameter) < 2 {
		return nil,fmt.Errorf("Error: Unrecognized mode parameter %s", restOfParameter)
	}
	
	switch restOfParameter[0:2] {
		case "IN":
			moParameter.units = UNITS_IN
		
		case "MM":
			moParameter.units = UNITS_MM
		
		default:
			return nil,fmt.Errorf("Error: Unrecognized mode parameter %s", restOfParameter)
	}
	
	// If we're here, we've successfully parsed the MO parameter, so update the parse environment
	env.unitsSet = true
	
	return moParameter,nil
}

func parseADParameter(adParameter *ApertureDefinitionParameter, restOfParameter string, env *ParseEnvironment) (DataBlock, error) {
	parsedAD := adParameterRegex.FindAllStringSubmatch(restOfParameter, -1)
	
	// Make sure we captured the number of subexpressions we expected
	if len(parsedAD) != 1 {
		return nil,fmt.Errorf("Unable to parse AD Parameter %s: error 1", restOfParameter)
	} else if len(parsedAD[0]) != 4 {
		return nil,fmt.Errorf("Unable to parse AD Parameter %s: error 2", restOfParameter)
	}
	
	newADParam := new(ApertureDefinitionParameter)
	
	// Parse the D code
	if dCode,err := strconv.ParseInt(parsedAD[0][1], 10, 32); err != nil {
		return nil,err
	} else {
		if dCode < 10 {
			return nil,fmt.Errorf("Aperture definition D codes must be 10 or larger.  Received %d", dCode)
		} else {
			newADParam.apertureNumber = int(dCode)
		}
	}
	
	// Make sure that this parameter hasn't already been defined.  It's illegal to re-use the same D code
	if _,exists := env.aperturesDefined[newADParam.apertureNumber]; exists {
		return nil,fmt.Errorf("Illegal duplicate aperture D code encountered: %d", newADParam.apertureNumber)
	}
	
	// Parse the aperture type
	switch parsedAD[0][2] {
		case "C":
			return parseCircleAperture(newADParam, parsedAD[0][3], env)
		
		case "R":
			return parseRectangleAperture(newADParam, parsedAD[0][3], env)
		
		case "O":
			return parseObroundAperture(newADParam, parsedAD[0][3], env)
		
		case "P":
			return parsePolygonAperture(newADParam, parsedAD[0][3], env)
		
		default:
			// Macro apertures take both a name, which is in the "type" slot,
			// and potentially modifiers, which are in the same modifier slot
			return parseMacroAperture(newADParam, parsedAD[0][2], parsedAD[0][3], env)
	}

	panic("End of parseADParameter: Shouldn't be here")
}

func parseCircleAperture(adParameter *ApertureDefinitionParameter, modifiers string, env *ParseEnvironment) (DataBlock, error) {
	parsedModifiers := strings.FieldsFunc(modifiers, modifierFieldsFunc)
	
	if len(parsedModifiers) < 1 {
		return nil,fmt.Errorf("Circle aperture definition missing required diameter modifier.  Received %s", modifiers)
	}
	
	// Set the aperture type
	adParameter.apertureType = CIRCLE_APERTURE
	
	newAperture := new(CircleAperture)
	
	// Parse circle diameter
	if diameter,err := strconv.ParseFloat(parsedModifiers[0], 64); err != nil {
		return nil,err
	} else {
		if diameter < 0 {
			return nil,fmt.Errorf("Circle aperture diameter must be 0 or greater.  Receive %f", diameter)
		} else {
			newAperture.diameter = diameter
		}
	}
	
	// If the aperture has a hole, parse it
	if len(parsedModifiers) > 1 {
		// If the hole parses correctly, parseApertureHole will add the parsed hole to the aperture struct
		if err := parseApertureHole(newAperture, parsedModifiers[1:]); err != nil {
			return nil,err
		}
	}
	
	// Save the aperture number
	newAperture.apertureNumber = adParameter.apertureNumber
	
	adParameter.aperture = newAperture
	
	// If we're here, we've successfully parsed the new aperture, so update the parse environment
	env.aperturesDefined[adParameter.apertureNumber] = true

	return adParameter,nil
}

func parseRectangleAperture(adParameter *ApertureDefinitionParameter, modifiers string, env *ParseEnvironment) (DataBlock, error) {
	parsedModifiers := strings.FieldsFunc(modifiers, modifierFieldsFunc)
	
	if len(parsedModifiers) < 2 {
		return nil,fmt.Errorf("Rectangle aperture definition missing required x size or y size modifier.  Received %s", modifiers)
	}
	
	// Set the aperture type
	adParameter.apertureType = RECTANGLE_APERTURE
	
	newAperture := new(RectangleAperture)
	
	// Parse rectangle x size
	if xSize,err := strconv.ParseFloat(parsedModifiers[0], 64); err != nil {
		return nil,err
	} else {
		if xSize <= 0 {
			return nil,fmt.Errorf("Rectangle aperture x size must be greater than 0.  Received %f", xSize)
		}
		newAperture.xSize = xSize
	}
	
	// Parse rectangle y size
	if ySize,err := strconv.ParseFloat(parsedModifiers[1], 64); err != nil {
		return nil,err
	} else {
		if ySize <= 0 {
			return nil,fmt.Errorf("Rectangle aperture y size must be greater than 0.  Received %f", ySize)
		}
		newAperture.ySize = ySize
	}
	
	// If the aperture has a hole, parse it
	if len(parsedModifiers) > 2 {
		// If the hole parses correctly, parseApertureHole will add the parsed hole to the aperture struct
		if err := parseApertureHole(newAperture, parsedModifiers[2:]); err != nil {
			return nil,err
		}
	}
	
	// Save the aperture number
	newAperture.apertureNumber = adParameter.apertureNumber
	
	adParameter.aperture = newAperture
	
	// If we're here, we've successfully parsed the new aperture, so update the parse environment
	env.aperturesDefined[adParameter.apertureNumber] = true

	return adParameter,nil
}

func parseObroundAperture(adParameter *ApertureDefinitionParameter, modifiers string, env *ParseEnvironment) (DataBlock, error) {
	parsedModifiers := strings.FieldsFunc(modifiers, modifierFieldsFunc)
	
	if len(parsedModifiers) < 2 {
		return nil,fmt.Errorf("Obround aperture definition missing required x size or y size modifier.  Received %s", modifiers)
	}
	
	// Set the aperture type
	adParameter.apertureType = OBROUND_APERTURE
	
	newAperture := new(ObroundAperture)
	
	// Parse obround x size
	if xSize,err := strconv.ParseFloat(parsedModifiers[0], 64); err != nil {
		return nil,err
	} else {
		if xSize <= 0 {
			return nil,fmt.Errorf("Obround aperture x size must be greater than 0.  Received %f", xSize)
		}
		newAperture.xSize = xSize
	}
	
	// Parse obround y size
	if ySize,err := strconv.ParseFloat(parsedModifiers[1], 64); err != nil {
		return nil,err
	} else {
		if ySize <= 0 {
			return nil,fmt.Errorf("Obround aperture y size must be greater than 0.  Received %f", ySize)
		}
		newAperture.ySize = ySize
	}
	
	// If the aperture has a hole, parse it
	if len(parsedModifiers) > 2 {
		// If the hole parses correctly, parseApertureHole will add the parsed hole to the aperture struct
		if err := parseApertureHole(newAperture, parsedModifiers[2:]); err != nil {
			return nil,err
		}
	}
	
	// Save the aperture number
	newAperture.apertureNumber = adParameter.apertureNumber
	
	adParameter.aperture = newAperture
	
	// If we're here, we've successfully parsed the new aperture, so update the parse environment
	env.aperturesDefined[adParameter.apertureNumber] = true

	return adParameter,nil
}

func parsePolygonAperture(adParameter *ApertureDefinitionParameter, modifiers string, env *ParseEnvironment) (DataBlock, error) {
	parsedModifiers := strings.FieldsFunc(modifiers, modifierFieldsFunc)
	
	if len(parsedModifiers) < 2 {
		return nil,fmt.Errorf("Polygon aperture definition missing required diameter or number of vertices modifier.  Received %s", modifiers)
	}
	
	// Set the aperture type
	adParameter.apertureType = POLYGON_APERTURE
	
	newAperture := new(PolygonAperture)
	
	// Parse polygon diameter
	if diameter,err := strconv.ParseFloat(parsedModifiers[0], 64); err != nil {
		return nil,err
	} else {
		if diameter <= 0 {
			return nil,fmt.Errorf("Polygon aperture diameter must be greater than 0.  Received %f", diameter)
		}
		newAperture.outerDiameter = diameter
	}
	
	// Parse polygon num vertices
	if vertices,err := strconv.ParseInt(parsedModifiers[1], 10, 32); err != nil {
		return nil,err
	} else {
		if (vertices < 3) || (vertices > 12) {
			return nil,fmt.Errorf("Polygon aperture number of vertices must be between 3 and 12.  Received %d", vertices)
		}
		newAperture.numVertices = int(vertices)
	}
	
	// Parse polygon rotation, if provided
	if len(parsedModifiers) > 2 {
		if rotation,err := strconv.ParseFloat(parsedModifiers[2], 64); err != nil {
			return nil,err
		} else {
			newAperture.rotationDegrees = rotation
		}
	} else {
		newAperture.rotationDegrees = 0.0
	}
	
	// If the aperture has a hole, parse it
	if len(parsedModifiers) > 3 {
		// If the hole parses correctly, parseApertureHole will add the parsed hole to the aperture struct
		if err := parseApertureHole(newAperture, parsedModifiers[3:]); err != nil {
			return nil,err
		}
	}
	
	// Save the aperture number
	newAperture.apertureNumber = adParameter.apertureNumber
	
	adParameter.aperture = newAperture
	
	// If we're here, we've successfully parsed the new aperture, so update the parse environment
	env.aperturesDefined[adParameter.apertureNumber] = true

	return adParameter,nil
}

func parseMacroAperture(adParameter *ApertureDefinitionParameter, name string, modifiers string, env *ParseEnvironment) (DataBlock, error) {
	adParameter.apertureType = MACRO_APERTURE
	adParameter.aperture = &MacroAperture{adParameter.apertureNumber, name}
	
	// If we're here, we've successfully parsed the new aperture, so update the parse environment
	env.aperturesDefined[adParameter.apertureNumber] = true
	
	return adParameter,nil
}

func parseApertureHole(aperture Aperture, holeModifiers []string) error {
	switch len(holeModifiers) {
		case 1:
			if diameter,err := strconv.ParseFloat(holeModifiers[0], 64); err != nil {
				return err
			} else {
				if diameter < 0 {
					return fmt.Errorf("Aperture circular hole must have diameter >= 0.  Received %f", diameter)
				} else {
					aperture.SetHole(&CircularHole{diameter})
				}
			}
			
		case 2:
			if xSize,err := strconv.ParseFloat(holeModifiers[0], 64); err != nil {
				return err
			} else {
				if ySize,err := strconv.ParseFloat(holeModifiers[1], 64); err != nil {
					return err
				} else {
					if xSize < 0 {
						return fmt.Errorf("Aperture square hole must have x size >= 0.  Received %f", xSize)
					}
					
					if ySize < 0 {
						return fmt.Errorf("Aperture square hole must have y size >= 0.  Received %f", ySize)
					}
					
					aperture.SetHole(&RectangularHole{xSize, ySize})
				}
			}
		
		default:
			return fmt.Errorf("Unexpected number of hole modifiers for aperture.  Must be 1 or 2, received %d", len(holeModifiers))
	}
	
	return nil
}

func parseAMParameter(amParameter *ApertureMacroParameter, restOfParameter string) (DataBlock, error) {
	// First, split the various data blocks apart
	blocks := strings.Split(restOfParameter, "*")
	
	fmt.Printf("Split AM blocks: %v\n", blocks)
	
	// The aperture macro parameter must have at least one block (the name)
	if len(blocks) < 1 {
		return nil,fmt.Errorf("Aperture Macro must have at least one data block")
	}
	
	// Populate the name
	amParameter.macroName = blocks[0]
	
	// Parse the rest of the macro
	return parseApertureMacro(amParameter, blocks[1:])
}

func parseSRParameter(srParameter *StepAndRepeatParameter, restOfParameter string) (DataBlock, error) {
	parsedSR := srParameterRegex.FindAllStringSubmatch(restOfParameter, -1)
	
	// Make sure we captured the number of subexpressions we expected
	if len(parsedSR) != 1 {
		return nil,fmt.Errorf("Unable to parse SR Parameter %s: error 1", restOfParameter)
	} else if len(parsedSR[0]) != 5 {
		return nil,fmt.Errorf("Unable to parse SR Parameter %s: error 2", restOfParameter)
	}
	
	// Parse X repeats
	if len(parsedSR[0][1]) > 0 {
		if xRepeats,err := strconv.ParseInt(parsedSR[0][1], 10, 32); err != nil {
			return nil,err
		} else {
			if xRepeats <= 0 {
				return nil,fmt.Errorf("X repeats must be a strictly positive integer.  Received %d", xRepeats)
			} else {
				srParameter.xRepeats = int(xRepeats)
			}
		}
	} else {
		srParameter.xRepeats = 1
	}

	// Parse Y repeats
	if len(parsedSR[0][2]) > 0 {
		if yRepeats,err := strconv.ParseInt(parsedSR[0][2], 10, 32); err != nil {
			return nil,err
		} else {
			if yRepeats <= 0 {
				return nil,fmt.Errorf("Y repeats must be a strictly positive integer.  Received %d", yRepeats)
			} else {
				srParameter.yRepeats = int(yRepeats)
			}
		}
	} else {
		srParameter.yRepeats = 1
	}
	
	// Parse I step
	if len(parsedSR[0][3]) > 0 {
		if iStep,err := strconv.ParseFloat(parsedSR[0][3], 64); err != nil {
			return nil,err
		} else {
			if iStep < 0 {
				return nil,fmt.Errorf("I Step must be non-negative decimal number.  Received %f", iStep)
			} else {
				srParameter.xStepDistance = iStep
			}
		}
	} else {
		// I step must be provided if X repeats is >1.  Otherwise, it's fine to leave it at
		// its default 0 value
		if srParameter.xRepeats > 1 {
			return nil,fmt.Errorf("If X repeats is >1, I step must be provided")
		}
	}
	
	// Parse J step
	if len(parsedSR[0][4]) > 0 {
		if jStep,err := strconv.ParseFloat(parsedSR[0][4], 64); err != nil {
			return nil,err
		} else {
			if jStep < 0 {
				return nil,fmt.Errorf("J Step must be non-negative decimal number.  Received %f", jStep)
			} else {
				srParameter.yStepDistance = jStep
			}
		}
	} else {
		// J step must be provided if Y repeats is >1.  Otherwise, it's fine to leave it at
		// its default 0 value
		if srParameter.yRepeats > 1 {
			return nil,fmt.Errorf("If Y repeats is >1, J step must be provided")
		}
	}

	return srParameter,nil
}

func parseLPParameter(lpParameter *LevelPolarityParameter, restOfParameter string) (DataBlock, error) {

	// Don't need to check for string length here, because we're only looking for one character
	// and the previous parsing steps have already determined that there is at least one character left
	
	switch restOfParameter[0] {
		case 'C':
			lpParameter.polarity = CLEAR_POLARITY
		
		case 'D':
			lpParameter.polarity = DARK_POLARITY
		
		default:
			return nil,fmt.Errorf("Unknown level polarity argument: %c", restOfParameter[0])
	}
	
	return lpParameter,nil
}

func modifierFieldsFunc(char rune) bool {
	return char == 'X'
}