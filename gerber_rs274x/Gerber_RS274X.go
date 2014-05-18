package gerber_rs274x

import (
	"io"
	"bufio"
	"fmt"
	"regexp"
	cairo "github.com/ungerik/go-cairo"
)

var coordDataBlockRegex *regexp.Regexp
var parameterOrDataBlockRegex *regexp.Regexp
var dataBlockRegex *regexp.Regexp
var dCodeDataBlockRegex *regexp.Regexp
var coordinateDataBlockRegex *regexp.Regexp
var fsParameterRegex *regexp.Regexp
var srParameterRegex *regexp.Regexp
var adParameterRegex *regexp.Regexp
var amVariableDefinitionRegex *regexp.Regexp

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

type ScalingParms struct {
	scaleFactor float64
	xOffset float64
	yOffset float64
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
	
	amVariableDefinitionRegex = regexp.MustCompile(`\$(?P<varNum>[[:digit:]]+)=(?P<varExp>[[:digit:]$.+-x/]+)`)
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

func GenerateSurface(outFileName string, parsedFile []DataBlock) error {
	
	width := 800
	height := 800
	
	// First, need to do a full render of the file, just keeping track of the bounds
	// of the generated image, so we can do the proper scaling when we render it for real
	gfxStateBounds := newGraphicsState(nil, 0, 0)
	bounds := newImageBounds()
	
	for _,dataBlock := range parsedFile {
		if err := dataBlock.ProcessDataBlockBoundsCheck(bounds, gfxStateBounds); err != nil {
			return err
		}
	}
	
	fmt.Printf("X Bounds: (%f %f) Y Bounds: (%f %f)\n", bounds.xMin, bounds.xMax, bounds.yMin, bounds.yMax)
	
	// Set up the graphics state for the actual drawing
	gfxState := newGraphicsState(bounds, width, height)
	
	// Construct the surface we're drawing to
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, width, height)
	surface.SetAntialias(cairo.ANTIALIAS_NONE)
	
	// This is important for regions with cut-ins.  If we leave the fill rule the default (winding),
	// cut-ins don't render correctly
	surface.SetFillRule(cairo.FILL_RULE_EVEN_ODD)
	// Invert the Y-axis.  This is to correct for the difference in coordinate frames between the gerber file and cairo
	surface.Scale(1.0, -1.0)
	surface.Translate(0.0, float64(-height))
	// Apply the x and y offsets as translations to the surface
	surface.Translate(gfxState.xOffset, gfxState.yOffset)
	
	//TODO: For testing
	surface.Save()
	surface.Scale(gfxState.scaleFactor, gfxState.scaleFactor)
	
	for _,dataBlock := range parsedFile {
		if err := dataBlock.ProcessDataBlockSurface(surface, gfxState); err != nil {
			gfxState.releaseRenderedSurfaces()
			surface.Finish()
			return err
		}
	}
	gfxState.releaseRenderedSurfaces()
	
	surface.WriteToPNG(outFileName)
	surface.Finish()
	
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
