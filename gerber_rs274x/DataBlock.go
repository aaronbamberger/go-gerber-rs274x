package gerber_rs274x

import "github.com/ajstarks/svgo"

type DataBlock interface {
	DataBlockPlaceholder()
	ProcessDataBlockSVG(svg *svg.SVG, gfxState *GraphicsState) error
}