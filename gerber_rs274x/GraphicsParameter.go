package gerber_rs274x

import cairo "github.com/ungerik/go-cairo"

type ImageNameParameter struct {
	paramCode ParameterCode
	name string
}

type ImageRotationParameter struct {
	paramCode ParameterCode
	rotation int
}

type OffsetParameter struct {
	paramCode ParameterCode
	axisAOffset float64
	axisBOffset float64
}

type AxisSelectParameter struct {
	paramCode ParameterCode
	isAXBY bool
}

type ImagePolarityParameter struct {
	paramCode ParameterCode
	polarity Polarity
}

type ScaleFactorParameter struct {
	paramCode ParameterCode
	axisAScale float64
	axisBScale float64
}

type LevelNameParameter struct {
	paramCode ParameterCode
	name string
}

type MirrorImageParameter struct {
	paramCode ParameterCode
	axisAMirror bool
	axisBMirror bool
}

func (apertureMacro *ApertureMacroParameter) DataBlockPlaceholder() {

}

func (apertureMacro *ApertureMacroParameter) ProcessDataBlockBoundsCheck(imageBounds *ImageBounds, gfxState *GraphicsState) error {
	//TODO: Implement this
	return nil
}

func (apertureMacro *ApertureMacroParameter) ProcessDataBlockSurface(surface *cairo.Surface, gfxState *GraphicsState) error {
	//TODO: Implement this
	return nil
}

func (imageName *ImageNameParameter) DataBlockPlaceholder() {

}

func (imageRotation *ImageRotationParameter) DataBlockPlaceholder() {

}

func (offset *OffsetParameter) DataBlockPlaceholder() {

}

func (axisSelect *AxisSelectParameter) DataBlockPlaceholder() {

}

func (imagePolarity *ImagePolarityParameter) DataBlockPlaceholder() {

}

func (scaleFactor *ScaleFactorParameter) DataBlockPlaceholder() {

}

func (levelName *LevelNameParameter) DataBlockPlaceholder() {

}

func (mirrorImage *MirrorImageParameter) DataBlockPlaceholder() {

}
