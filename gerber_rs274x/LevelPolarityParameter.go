package gerber_rs274x

import (
	"fmt"
	cairo "github.com/ungerik/go-cairo"
)

type LevelPolarityParameter struct {
	paramCode ParameterCode
	polarity Polarity
}

func (levelPolarity *LevelPolarityParameter) DataBlockPlaceholder() {

}

func (levelPolarity *LevelPolarityParameter) ProcessDataBlockBoundsCheck(imageBounds *ImageBounds, gfxState *GraphicsState) error {
	gfxState.currentLevelPolarity = levelPolarity.polarity
	
	return nil
}

func (levelPolarity *LevelPolarityParameter) ProcessDataBlockSurface(surface *cairo.Surface, gfxState *GraphicsState) error {
	gfxState.currentLevelPolarity = levelPolarity.polarity
	
	return nil
}

func (lpParam *LevelPolarityParameter) String() string {
	var levelPolarity string
	
	switch lpParam.polarity {
		case CLEAR_POLARITY:
			levelPolarity = "Clear"
			
		case DARK_POLARITY:
			levelPolarity = "Dark"
			
		default:
			levelPolarity = "Unknown"
	}
	
	return fmt.Sprintf("{LP, Polarity: %s}", levelPolarity)
}