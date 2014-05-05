package gerber_rs274x

import (
	"fmt"
	cairo "github.com/ungerik/go-cairo"
)

type IgnoreDataBlock struct {
	comment string
}

func (ignoreDataBlock *IgnoreDataBlock) DataBlockPlaceholder() {

}

func (ignoreDataBlock *IgnoreDataBlock) ProcessDataBlockBoundsCheck(imageBounds *ImageBounds, gfxState *GraphicsState) error {
	// This is a comment, so it doesn't change the graphics state or draw anything
	return nil
}

func (ignoreDataBlock *IgnoreDataBlock) ProcessDataBlockSurface(surface *cairo.Surface, gfxState *GraphicsState) error {
	// This is a comment, so it doesn't change the graphics state or draw anything
	return nil
}

func (ignoreDataBlock *IgnoreDataBlock) String() string {
	return fmt.Sprintf("{COMMENT, %s}", ignoreDataBlock.comment)
}