package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

type StepAndRepeatParameter struct {
	paramCode ParameterCode
	xRepeats int
	yRepeats int
	xStepDistance float64
	yStepDistance float64
}

func (stepAndRepeat *StepAndRepeatParameter) DataBlockPlaceholder() {

}

func (stepAndRepeat *StepAndRepeatParameter) ProcessDataBlockSVG(svg *svg.SVG, gfxState *GraphicsState) error {
	//TODO: Implement this
	return nil
}

func (srParam *StepAndRepeatParameter) String() string {
	return fmt.Sprintf("{SR, X Repeats: %d, Y Repeats: %d, I Step: %f, J Step: %f}", srParam.xRepeats, srParam.yRepeats, srParam.xStepDistance, srParam.yStepDistance)
}

