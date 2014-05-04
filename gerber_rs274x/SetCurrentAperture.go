package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

type SetCurrentAperture struct {
	apertureNumber int
}

func (setCurrentAperture *SetCurrentAperture) DataBlockPlaceholder() {

}

func (setCurrentAperture *SetCurrentAperture) ProcessDataBlockSVG(svg *svg.SVG, gfxState *GraphicsState) error {
	// Make sure the aperture we're trying to switch to has already been defined
	if _,exists := gfxState.apertures[setCurrentAperture.apertureNumber]; !exists {
		return fmt.Errorf("Unable to switch to undefined aperture %d", setCurrentAperture.apertureNumber)
	}

	gfxState.currentAperture = setCurrentAperture.apertureNumber
	gfxState.apertureSet = true
	
	return nil
}

func (setCurrentAperture *SetCurrentAperture) String() string {
	return fmt.Sprintf("{SET APERTURE, Aperture: %d}", setCurrentAperture.apertureNumber)
}