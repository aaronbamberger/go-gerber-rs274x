package gerber_rs274x

import (
	"fmt"
	"math"
	cairo "github.com/ungerik/go-cairo"
)

type MacroAperture struct {
	apertureNumber int
	macroName string
	env *ExpressionEnvironment
	xSize float64
	ySize float64
	sizeCalculated bool
}

func (aperture *MacroAperture) AperturePlaceholder() {

}

func (aperture *MacroAperture) GetHole() Hole {
	return nil
}

func (aperture *MacroAperture) SetHole(hole Hole) {
	
}

func (aperture *MacroAperture) GetMinSize() float64 {
	//TODO: Implement appropriately
	return math.MaxFloat64
}

func (aperture *MacroAperture) DrawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error {
	return nil
}

func (aperture *MacroAperture) DrawApertureSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	return nil
}

func (aperture *MacroAperture) calculateApertureSize(macroDataBlocks []ApertureMacroDataBlock) {
	// We need to execute the entire macro to calculate the size, and this will pollute the enviroment
	// for when we want to actually render the aperture, so we need to create a copy of the environment to use
	// while calculating size
	sizeEnv := NewExpressionEnvironment()
	for key,value := range aperture.env.variables {
		sizeEnv.setVariableValue(key,value)
	}
	
	for _,dataBlock := range macroDataBlocks {
		switch dataBlockValue := dataBlock.(type) {
			case *ApertureMacroComment:
				// Nothing to do here
			
			case *ApertureMacroVariableDefinition:
				// Need to update the expression environment
				sizeEnv.setVariableValue(dataBlockValue.variableNumber, dataBlockValue.value.EvaluateExpression(sizeEnv))
			
			case AperturePrimitive:
				
				
		}	
	}
	
}

func (aperture *MacroAperture) String() string {
	return fmt.Sprintf("{MA, Name: %s}", aperture.macroName)
}