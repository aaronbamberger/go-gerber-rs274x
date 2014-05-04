package gerber_rs274x

import (
	"io"
	"bufio"
	"fmt"
	"regexp"
	"github.com/ajstarks/svgo"
)

var coordDataBlockRegex *regexp.Regexp
var parameterOrDataBlockRegex *regexp.Regexp
var dataBlockRegex *regexp.Regexp
var dCodeDataBlockRegex *regexp.Regexp
var coordinateDataBlockRegex *regexp.Regexp
var fsParameterRegex *regexp.Regexp
var srParameterRegex *regexp.Regexp
var adParameterRegex *regexp.Regexp

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

type ParseEnvironment struct {
	coordFormat CoordinateFormat
	unitsSet bool
	aperturesDefined map[int]bool
}

type GraphicsState struct {
	currentAperture int
	currentQuadrantMode FunctionCode
	currentInterpolationMode FunctionCode
	currentX float64
	currentY float64
	currentLevelPolarity Polarity
	regionModeOn bool
	xImageSize int
	yImageSize int
	fileComplete bool
	coordinateNotation CoordinateNotation
	drawPrecision float64
	
	// As we encounter aperture definitions, we save them
	// for later use while drawing
	apertures map[int]Aperture
	
	// Some of these default to undefined,
	// so we also need to keep track of when they get defined
	apertureSet bool
	quadrantModeSet bool
	interpolationModeSet bool
	coordinateNotationSet bool
}

func (gfxState *GraphicsState) updateCurrentCoordinate(newX float64, newY float64) {
	gfxState.currentX = newX
	gfxState.currentY = newY
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
	
	fsParameterRegex = regexp.MustCompile(`(?P<zeroOmissionMode>L|T)(?P<coordinateNotation>A|I)X(?P<xIntPositions>[[:digit:]]{1})(?P<xDecPositions>[[:digit:]]{1})Y(?P<yIntPositions>[[:digit:]]{1})(?P<yDecPositions>[[:digit:]]{1})`)
	
	srParameterRegex = regexp.MustCompile(`(?:X(?P<xRepeat>[[:digit:]]+))?(?:Y(?P<yRepeat>[[:digit:]]+))?(?:I(?P<iStep>[[:digit:]]+\.?[[:digit:]]*))?(?:J(?P<jStep>[[:digit:]]+\.?[[:digit:]]*))?`)
	
	adParameterRegex = regexp.MustCompile(`D(?P<dCode>[[:digit:]]*)(?P<apertureType>[[:alnum:]_\+\-/\!\?<>"'\(\){}\.\\\|\&@# ]+),?(?P<modifiers>[[:digit:]\.X]*)`)
}

func ParseGerberFile(in io.Reader) (parsedFile []DataBlock, err error) {
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
	parseEnv := newParseEnv()
	parsedFile = make([]DataBlock, 0, 100)
	
	for index,submatch := range results {
		if len(submatch) != 3 {
			return nil,fmt.Errorf("Error (token %d): Parse error on command %v\n", index, submatch)
		}
		
		if len(submatch[1]) > 0 {
			// Parsing Parameter
			fmt.Printf("Token %d, Parsed parameter: %s\n", index, submatch[1])
			if parameter,err := parseParameter(submatch[1], parseEnv); err != nil {
				fmt.Printf("Parse error for parameter %s: %s\n", submatch[1], err.Error())
			} else {
				parsedFile = append(parsedFile, parameter)
			}
		} else if len(submatch[2]) > 0 {
			// Parsing non-parameter data block
			if dataBlock,err := parseDataBlock(submatch[2], parseEnv); err != nil {
				fmt.Printf("Parse Error for block %s: %s\n", submatch[2], err.Error())
			} else {
				parsedFile = append(parsedFile, dataBlock)
			}
		} else {
			return nil,fmt.Errorf("Error (token %d): Not parameter or data block: %v\n", index, submatch)
		}
	}
	
	for index,dataBlock := range parsedFile {
		fmt.Printf("Parsed data block %3d: %v\n", index, dataBlock)
	}
	
	return parsedFile,nil
}

func GenerateSVG(out io.Writer, parsedFile []DataBlock) error {
	
	width := 4000
	height := 2000
	
	// Set up the initial graphics state
	gfxState := newGraphicsState(width, height)
	
	canvas := svg.New(out)
	canvas.Start(width, height)
	
	for _,dataBlock := range parsedFile {
		if err := dataBlock.ProcessDataBlockSVG(canvas, gfxState); err != nil {
			return err
		}
	}
	
	canvas.End()
	
	// Make sure that the entire file was rendered
	if !gfxState.fileComplete {
		return fmt.Errorf("Render of file completed without reaching end of file code (M02)")
	}
	
	return nil
}

func newParseEnv() *ParseEnvironment {
	parseEnv := new(ParseEnvironment)
	parseEnv.aperturesDefined = make(map[int]bool, 10) // We'll start with an initial capacity of 10, it will grow as necessary
	
	return parseEnv
}

func newGraphicsState(xImageSize int, yImageSize int) *GraphicsState {
	graphicsState := new(GraphicsState)
	
	graphicsState.currentLevelPolarity = DARK_POLARITY
	graphicsState.xImageSize = xImageSize
	graphicsState.yImageSize = yImageSize
	graphicsState.apertures = make(map[int]Aperture, 10) // Start with an initial capacity of 10 apertures, will grow as needed
	
	// All other settings are fine with their go defaults
	// Current aperture: Doesn't matter since it's undefined by default
	// Current quadrant mode: Doesn't matter since it's undefined by default
	// Current interpolation mode: Doesn't matter since it's undefined by default
	// Coordinate notation: Doesn't matter since it's undefined by default
	// Current x: 0 is correct
	// Current y: 0 is correct
	// Region mode on: false is correct
	// Aperture set: false is correct
	// Quadrant mode set: false is correct
	// Interpolation mode set: false is correct
	// Region mode on: false is correct
	// File complete: false is correct
	// Coordinate notation set: false is correct
	
	return graphicsState 
}
