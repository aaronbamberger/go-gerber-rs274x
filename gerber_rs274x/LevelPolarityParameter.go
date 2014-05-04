package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

type LevelPolarityParameter struct {
	paramCode ParameterCode
	polarity Polarity
}

func (levelPolarity *LevelPolarityParameter) DataBlockPlaceholder() {

}

func (levelPolarity *LevelPolarityParameter) ProcessDataBlockSVG(svg *svg.SVG, gfxState *GraphicsState) error {
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