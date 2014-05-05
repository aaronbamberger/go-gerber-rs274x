package gerber_rs274x

import cairo "github.com/ungerik/go-cairo"

type DataBlock interface {
	DataBlockPlaceholder()
	ProcessDataBlockBoundsCheck(imageBounds *ImageBounds, gfxState *GraphicsState) error
	ProcessDataBlockSurface(surface *cairo.Surface, gfxState *GraphicsState) error
}