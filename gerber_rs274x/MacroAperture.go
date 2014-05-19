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
	xMin float64
	xMax float64
	yMin float64
	yMax float64
	boundsCalculated bool
}

func (aperture *MacroAperture) AperturePlaceholder() {

}

func (aperture *MacroAperture) GetApertureNumber() int {
	return aperture.apertureNumber
}

func (aperture *MacroAperture) GetHole() Hole {
	return nil
}

func (aperture *MacroAperture) SetHole(hole Hole) {
	
}

func (aperture *MacroAperture) GetMinSize(gfxState *GraphicsState) float64 {
	if !aperture.boundsCalculated {
		// If the bounds haven't been calculated yet, do it now
		// First, retrieve the aperture macro from the graphics state
		if macro,found := gfxState.apertureMacros[aperture.macroName]; !found {
			//TODO: Figure out better error behavior for this
			return math.MaxFloat64
		} else {
			aperture.calculateApertureSize(macro)
		}
	}
	
	return math.Min(aperture.xMax - aperture.xMin, aperture.yMax - aperture.yMin)
}

func (aperture *MacroAperture) DrawApertureBoundsCheck(bounds *ImageBounds, gfxState *GraphicsState, x float64, y float64) error {
	if !aperture.boundsCalculated {
		// If the bounds haven't been calculated yet, do it now
		// First, retrieve the aperture macro from the graphics state
		if macro,found := gfxState.apertureMacros[aperture.macroName]; !found {
			return fmt.Errorf("Attempt to assign aperture %s to D code %d before it has been defined", aperture.macroName, aperture.apertureNumber)
		} else {
			aperture.calculateApertureSize(macro)
		}
	}
	
	xMin := x - aperture.xMin
	xMax := x + aperture.xMax
	yMin := y - aperture.yMin
	yMax := y + aperture.yMax
	
	bounds.updateBounds(xMin, xMax, yMin, yMax)
	
	return nil
}

func (aperture *MacroAperture) DrawApertureSurface(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	radiusX := (aperture.xMax - aperture.xMin) / 2.0
	radiusY := (aperture.yMax - aperture.yMin) / 2.0
	
	correctedX := x - radiusX
	correctedY := y - radiusY
	
	return renderApertureToSurface(aperture, surface, gfxState, correctedX, correctedY)
}

func (aperture *MacroAperture) DrawApertureSurfaceNoHole(surface *cairo.Surface, gfxState *GraphicsState, x float64, y float64) error {
	radiusX := (aperture.xMax - aperture.xMin) / 2.0
	radiusY := (aperture.yMax - aperture.yMin) / 2.0
	
	correctedX := x - radiusX
	correctedY := y - radiusY
	
	return renderApertureNoHoleToSurface(aperture, surface, gfxState, correctedX, correctedY)
}

func (aperture *MacroAperture) StrokeApertureLinear(surface *cairo.Surface, gfxState *GraphicsState, startX float64, startY float64, endX float64, endY float64) error {
	return nil
}

func (aperture *MacroAperture) StrokeApertureClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *MacroAperture) StrokeApertureCounterClockwise(surface *cairo.Surface, gfxState *GraphicsState, centerX float64, centerY float64, radius float64, startAngle float64, endAngle float64) error {
	return nil
}

func (aperture *MacroAperture) renderApertureToGraphicsState(gfxState *GraphicsState) {
	// This will render the aperture to a cairo surface the first time it is needed, then
	// cache it in the graphics state.  Subsequent draws of the aperture will used the cached surface
	
	rangeX := aperture.xMax - aperture.xMin
	rangeY := aperture.yMax - aperture.yMin
	radiusX := rangeX / 2.0
	radiusY := rangeY / 2.0
	
	// Construct the surface we're drawing to
	imageSizeX := int(math.Ceil(rangeX * gfxState.scaleFactor))
	imageSizeY := int(math.Ceil(rangeY * gfxState.scaleFactor))
	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, imageSizeX, imageSizeY)
	// Scale the surface so we can use unscaled coordinates in the primitive rendering routines
	surface.Scale(gfxState.scaleFactor, gfxState.scaleFactor)
	// Apply an offset to the surface, so that the origin is actually in the center of the image
	surface.Translate(radiusX, radiusY)
	
	// Set fill rule to Even/Odd so that rings render correctly
	surface.SetFillRule(cairo.FILL_RULE_EVEN_ODD)
	
	// Draw the aperture
	if gfxState.currentLevelPolarity == DARK_POLARITY {
		surface.SetSourceRGBA(0.0, 0.0, 0.0, 1.0)
	} else {
		surface.SetSourceRGBA(1.0, 1.0, 1.0, 1.0)
	}
	
	// Retrieve the macro from the graphics state
	if macro,found := gfxState.apertureMacros[aperture.macroName]; !found {
		//TODO: Figure out the error behavior, just print a warning for now
		fmt.Printf("Error: Attempt to render macro aperture %s before it has been defined\n", aperture.macroName)
	} else {
		for _,dataBlock := range macro {
			switch dataBlockValue := dataBlock.(type) {
				case *ApertureMacroComment:
					//Nothing to do here
					
				case *ApertureMacroVariableDefinition:
					// Need to update the expression environment
					aperture.env.setVariableValue(dataBlockValue.variableNumber, dataBlockValue.value.EvaluateExpression(aperture.env))
					
				case AperturePrimitive:
					if err := dataBlockValue.DrawPrimitiveToSurface(surface, aperture.env); err != nil {
						// TODO: Figure out the error behavior, just print a warning for now
						fmt.Printf("Error while attempting to render primitive on macro aperture %s: %s\n", aperture.macroName, err.Error())
					}		
			}
		}
	}
	
	surface.WriteToPNG(fmt.Sprintf("Aperture-%d.png", aperture.apertureNumber))
	
	gfxState.renderedApertures[aperture.apertureNumber] = surface
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
				xMin,xMax,yMin,yMax := dataBlockValue.GetPrimitiveBounds(sizeEnv)
				if xMin < aperture.xMin {
					aperture.xMin = xMin
				}
				if xMax > aperture.xMax {
					aperture.xMax = xMax
				}
				if yMin < aperture.yMin {
					aperture.yMin = yMin
				}
				if yMax > aperture.yMax {
					aperture.yMax = yMax
				}
		}
	}
	aperture.boundsCalculated = true
}

func (aperture *MacroAperture) String() string {
	return fmt.Sprintf("{MA, Name: %s}", aperture.macroName)
}