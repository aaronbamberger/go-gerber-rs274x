package gerber_rs274x

import (
	"fmt"
	"strconv"
	"strings"
)

type ApertureMacroParameter struct {
	paramCode ParameterCode
	macroName string
	dataBlocks []ApertureMacroDataBlock
}

type ApertureMacroDataBlock interface {
	ApertureMacroDataBlockPlaceholder()
}

type AperturePrimitive interface {
	ApertureMacroDataBlock
	AperturePrimitivePlaceholder()
}

type ApertureMacroVariableDefinition struct {
	variableNumber int
	value ApertureMacroExpression
}

type ApertureMacroComment struct {
	comment string
}

func (variableDefinition *ApertureMacroVariableDefinition) ApertureMacroDataBlockPlaceholder() {

}

func (comment *ApertureMacroComment) ApertureMacroDataBlockPlaceholder() {

}

func parseApertureMacro(amParameter *ApertureMacroParameter, dataBlocks []string) (*ApertureMacroParameter, error) {
	// Create the data blocks slice with the appropriate capacity 
	amParameter.dataBlocks = make([]ApertureMacroDataBlock, 0, len(dataBlocks))

	for _,dataBlock := range dataBlocks {
		switch dataBlock[0] {
			case '0':
				// Slice off the first two characters of the comment (the "0" comment specifier, and the opening space),
				// and append the comment block to the slice of parsed data blocks
				amParameter.dataBlocks = append(amParameter.dataBlocks, &ApertureMacroComment{dataBlock[2:]})
			
			case '$':
				// Parse the variable assignment
				varParts := amVariableDefinitionRegex.FindAllStringSubmatch(dataBlock, -1)
				
				// Make sure the regexp parse worked as expected
				if len(varParts) != 1 || len(varParts[0]) != 3 {
					return nil,fmt.Errorf("Error parsing aperture macro variable assignment %s", dataBlock)
				}
				
				// Parse the variable number
				if varNum,err := strconv.ParseInt(varParts[0][1], 10, 32); err != nil {
					return nil,err
				} else {
					// Parse the variable expression
					if expr,err := parseExpression(varParts[0][2]); err != nil {
						return nil,err
					} else {
						// Now that we've successfully parsed the variable number and expression, add a new variable definition to the macro parameter
						amParameter.dataBlocks = append(amParameter.dataBlocks, &ApertureMacroVariableDefinition{int(varNum), expr})
					}
				}
			
			default:
				// This is an aperture primitive, so parse accordingly
				if primitive,err := parseAperturePrimitive(dataBlock); err != nil {
					return nil,err
				} else {
					amParameter.dataBlocks = append(amParameter.dataBlocks, primitive)
				}
		}
	}

	return amParameter,nil
}

func parseAperturePrimitive(primitiveDefinition string) (AperturePrimitive, error) {
	splitPrimitive := strings.Split(primitiveDefinition, ",")

	// The primitive must have at least 1 part (the primitive code)
	if len(splitPrimitive) < 1 {
		return nil,fmt.Errorf("Primitive definition %s missing primitive code", primitiveDefinition)
	}
	
	switch splitPrimitive[0] {
		case "1":
			return parseCirclePrimitive(splitPrimitive[1:])
			
		case "2","20":
			return parseVectorLinePrimitive(splitPrimitive[1:])
			
		case "21":
			return parseCenterLinePrimitive(splitPrimitive[1:])
			
		case "22":
			return parseLowerLeftLinePrimitive(splitPrimitive[1:])
			
		case "4":
			return parseOutlinePrimitive(splitPrimitive[1:])
			
		case "5":
			return parsePolygonPrimitive(splitPrimitive[1:])
			
		case "6":
			return parseMoirePrimitive(splitPrimitive[1:])
			
		case "7":
			return parseThermalPrimitive(splitPrimitive[1:])
			
		default:
			return nil,fmt.Errorf("Unrecognized aperture primitive code: %s", splitPrimitive[0])
	}
}

func parseCirclePrimitive(modifiers []string) (AperturePrimitive, error) {
	// Check the number of modifiers
	if len(modifiers) != 4 {
		return nil,fmt.Errorf("Wrong number of modifiers for circle primitive.  Expected 4, received %d", len(modifiers))
	}
	
	// Parse the modifiers
	var exposure ApertureMacroExpression
	var diameter ApertureMacroExpression
	var centerX ApertureMacroExpression
	var centerY ApertureMacroExpression
	var err error
	
	if exposure,err = parseExpression(modifiers[0]); err != nil {
		return nil,err
	}
	
	if diameter,err = parseExpression(modifiers[1]); err != nil {
		return nil,err
	}
	
	if centerX,err = parseExpression(modifiers[2]); err != nil {
		return nil,err
	}
	
	if centerY,err = parseExpression(modifiers[3]); err != nil {
		return nil,err
	}
	
	return &CirclePrimitive{exposure, diameter, centerX, centerY},nil
}

func parseVectorLinePrimitive(modifiers []string) (AperturePrimitive, error) {
	// Check the number of modifiers
	if len(modifiers) != 7 {
		return nil,fmt.Errorf("Wrong number of modifiers for vector line primitive.  Expected 7, received %d", len(modifiers))
	}

	var exposure ApertureMacroExpression
	var lineWidth ApertureMacroExpression
	var startX ApertureMacroExpression
	var startY ApertureMacroExpression
	var endX ApertureMacroExpression
	var endY ApertureMacroExpression
	var rotation ApertureMacroExpression
	var err error
	
	if exposure,err = parseExpression(modifiers[0]); err != nil {
		return nil,err
	}
	
	if lineWidth,err = parseExpression(modifiers[1]); err != nil {
		return nil,err
	}
	
	if startX,err = parseExpression(modifiers[2]); err != nil {
		return nil,err
	}
	
	if startY,err = parseExpression(modifiers[3]); err != nil {
		return nil,err
	}
	
	if endX,err = parseExpression(modifiers[4]); err != nil {
		return nil,err
	}
	
	if endY,err = parseExpression(modifiers[5]); err != nil {
		return nil,err
	}
	
	if rotation,err = parseExpression(modifiers[6]); err != nil {
		return nil,err
	}

	return &VectorLinePrimitive{exposure, lineWidth, startX, startY, endX, endY, rotation},nil
}

func parseCenterLinePrimitive(modifiers []string) (AperturePrimitive, error) {
	// Check the number of modifiers
	if len(modifiers) != 6 {
		return nil,fmt.Errorf("Wrong number of modifiers for center line primitive.  Expected 6, received %d", len(modifiers))
	}
	
	var exposure ApertureMacroExpression
	var width ApertureMacroExpression
	var height ApertureMacroExpression
	var centerX ApertureMacroExpression
	var centerY ApertureMacroExpression
	var rotation ApertureMacroExpression
	var err error
	
	if exposure,err = parseExpression(modifiers[0]); err != nil {
		return nil,err
	}
	
	if width,err = parseExpression(modifiers[1]); err != nil {
		return nil,err
	}
	
	if height,err = parseExpression(modifiers[2]); err != nil {
		return nil,err
	}
	
	if centerX,err = parseExpression(modifiers[3]); err != nil {
		return nil,err
	}
	
	if centerY,err = parseExpression(modifiers[4]); err != nil {
		return nil,err
	}
	
	if rotation,err = parseExpression(modifiers[5]); err != nil {
		return nil,err
	}
	
	return &CenterLinePrimitive{exposure, width, height, centerX, centerY, rotation},nil
}

func parseLowerLeftLinePrimitive(modifiers []string) (AperturePrimitive, error) {
	// Check the number of modifiers
	if len(modifiers) != 6 {
		return nil,fmt.Errorf("Wrong number of modifiers for lower left line primitive.  Expected 6, received %d", len(modifiers))
	}
	
	var exposure ApertureMacroExpression
	var width ApertureMacroExpression
	var height ApertureMacroExpression
	var lowerLeftX ApertureMacroExpression
	var lowerLeftY ApertureMacroExpression
	var rotation ApertureMacroExpression
	var err error
	
	if exposure,err = parseExpression(modifiers[0]); err != nil {
		return nil,err
	}
	
	if width,err = parseExpression(modifiers[1]); err != nil {
		return nil,err
	}
	
	if height,err = parseExpression(modifiers[2]); err != nil {
		return nil,err
	}
	
	if lowerLeftX,err = parseExpression(modifiers[3]); err != nil {
		return nil,err
	}
	
	if lowerLeftY,err = parseExpression(modifiers[4]); err != nil {
		return nil,err
	}
	
	if rotation,err = parseExpression(modifiers[5]); err != nil {
		return nil,err
	}

	return &LowerLeftLinePrimitive{exposure, width, height, lowerLeftX, lowerLeftY, rotation},nil
}

func parseOutlinePrimitive(modifiers []string) (AperturePrimitive, error) {
	// We can't check the exact number of modifiers, because it depends on the number of vertices, which can't be determined
	// until the entire file is parsed (because it might depend on an argument)
	// We do know, however, that there must be at least 7 modifiers, so we check for that
	// Check the number of modifiers
	if len(modifiers) < 7 {
		return nil,fmt.Errorf("Wrong number of modifiers for outline primitive.  Expected at least 7, received %d", len(modifiers))
	}
	
	// We also know there needs to be an odd number of modifiers, so check that as well (an even number would break the subsequent loop logic)
	if len(modifiers) % 2 == 0 {
		return nil,fmt.Errorf("There must be an odd number of modifiers for an outline primitive.  Received %d", len(modifiers))
	}
	
	outlinePrimitive := new(OutlinePrimitive)
	outlinePrimitive.subsequentX = make([]ApertureMacroExpression, 0, 10) // We'll start both with 10 elements to begin with
	outlinePrimitive.subsequentY = make([]ApertureMacroExpression, 0, 10)
	
	var exposure ApertureMacroExpression
	var nSubsequentPoints ApertureMacroExpression
	var startX ApertureMacroExpression
	var startY ApertureMacroExpression
	var rotation ApertureMacroExpression
	var subsequentCoord ApertureMacroExpression
	var err error
	
	if exposure,err = parseExpression(modifiers[0]); err != nil {
		return nil,err
	}
	
	if nSubsequentPoints,err = parseExpression(modifiers[1]); err != nil {
		return nil,err
	}
	
	if startX,err = parseExpression(modifiers[2]); err != nil {
		return nil,err
	}
	
	if startY,err = parseExpression(modifiers[3]); err != nil {
		return nil,err
	}
	
	for point := 4; point < len(modifiers) - 1; point += 2 {
		// Parse subsequent x coordinate
		if subsequentCoord,err = parseExpression(modifiers[point]); err != nil {
			return nil,err
		}
		outlinePrimitive.subsequentX = append(outlinePrimitive.subsequentX, subsequentCoord)
		
		// Parse subsequent y coordinate
		if subsequentCoord,err = parseExpression(modifiers[point + 1]); err != nil {
			return nil,err
		}
		outlinePrimitive.subsequentY = append(outlinePrimitive.subsequentY, subsequentCoord)
	}
	
	if rotation,err = parseExpression(modifiers[len(modifiers) - 1]); err != nil {
		return nil,err
	}
	
	outlinePrimitive.exposure = exposure
	outlinePrimitive.nPoints = nSubsequentPoints
	outlinePrimitive.startX = startX
	outlinePrimitive.startY = startY
	outlinePrimitive.rotationAngle = rotation

	return outlinePrimitive,nil
}

func parsePolygonPrimitive(modifiers []string) (AperturePrimitive, error) {
	// Check the number of modifiers
	if len(modifiers) != 6 {
		return nil,fmt.Errorf("Wrong number of modifiers for polygon primitive.  Expected 6, received %d", len(modifiers))
	}
	
	var exposure ApertureMacroExpression
	var numVertices ApertureMacroExpression
	var centerX ApertureMacroExpression
	var centerY ApertureMacroExpression
	var diameter ApertureMacroExpression
	var rotation ApertureMacroExpression
	var err error
	
	if exposure,err = parseExpression(modifiers[0]); err != nil {
		return nil,err
	}
	
	if numVertices,err = parseExpression(modifiers[1]); err != nil {
		return nil,err
	}
	
	if centerX,err = parseExpression(modifiers[2]); err != nil {
		return nil,err
	}
	
	if centerY,err = parseExpression(modifiers[3]); err != nil {
		return nil,err
	}
	
	if diameter,err = parseExpression(modifiers[4]); err != nil {
		return nil,err
	}
	
	if rotation,err = parseExpression(modifiers[5]); err != nil {
		return nil,err
	}

	return &PolygonPrimitive{exposure, numVertices, centerX, centerY, diameter, rotation},nil
}

func parseMoirePrimitive(modifiers []string) (AperturePrimitive, error) {
	// Check the number of modifiers
	if len(modifiers) != 9 {
		return nil,fmt.Errorf("Wrong number of modifiers for moire primitive.  Expected 9, received %d", len(modifiers))
	}
	
	var centerX ApertureMacroExpression
	var centerY ApertureMacroExpression
	var outerDiameter ApertureMacroExpression
	var ringThickness ApertureMacroExpression
	var ringGap ApertureMacroExpression
	var maxRings ApertureMacroExpression
	var crosshairThickness ApertureMacroExpression
	var crosshairLength ApertureMacroExpression
	var rotation ApertureMacroExpression
	var err error
	
	if centerX,err = parseExpression(modifiers[0]); err != nil {
		return nil,err
	}
	
	if centerY,err = parseExpression(modifiers[1]); err != nil {
		return nil,err
	}
	
	if outerDiameter,err = parseExpression(modifiers[2]); err != nil {
		return nil,err
	}
	
	if ringThickness,err = parseExpression(modifiers[3]); err != nil {
		return nil,err
	}
	
	if ringGap,err = parseExpression(modifiers[4]); err != nil {
		return nil,err
	}
	
	if maxRings,err = parseExpression(modifiers[5]); err != nil {
		return nil,err
	}
	
	if crosshairThickness,err = parseExpression(modifiers[6]); err != nil {
		return nil,err
	}
	
	if crosshairLength,err = parseExpression(modifiers[7]); err != nil {
		return nil,err
	}
	
	if rotation,err = parseExpression(modifiers[8]); err != nil {
		return nil,err
	}

	return &MoirePrimitive{centerX, centerY, outerDiameter, ringThickness, ringGap, maxRings, crosshairThickness, crosshairLength, rotation},nil
}

func parseThermalPrimitive(modifiers []string) (AperturePrimitive, error) {
	// Check the number of modifiers
	if len(modifiers) != 6 {
		return nil,fmt.Errorf("Wrong number of modifiers for thermal primitive.  Expected 6, received %d", len(modifiers))
	}
	
	var centerX ApertureMacroExpression
	var centerY ApertureMacroExpression
	var outerDiameter ApertureMacroExpression
	var innerDiameter ApertureMacroExpression
	var gapThickness ApertureMacroExpression
	var rotation ApertureMacroExpression
	var err error
	
	if centerX,err = parseExpression(modifiers[0]); err != nil {
		return nil,err
	}
	
	if centerY,err = parseExpression(modifiers[1]); err != nil {
		return nil,err
	}
	
	if outerDiameter,err = parseExpression(modifiers[2]); err != nil {
		return nil,err
	}
	
	if innerDiameter,err = parseExpression(modifiers[3]); err != nil {
		return nil,err
	}
	
	if gapThickness,err = parseExpression(modifiers[4]); err != nil {
		return nil,err
	}
	
	if rotation,err = parseExpression(modifiers[5]); err != nil {
		return nil,err
	}
	
	return &ThermalPrimitive{centerX, centerY, outerDiameter, innerDiameter, gapThickness, rotation},nil
}