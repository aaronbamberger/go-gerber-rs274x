package gerber_rs274x

import (
	"fmt"
	"github.com/ajstarks/svgo"
)

type IgnoreDataBlock struct {
	comment string
}

func (ignoreDataBlock *IgnoreDataBlock) DataBlockPlaceholder() {

}

func (ignoreDataBlock *IgnoreDataBlock) ProcessDataBlockSVG(svg *svg.SVG, gfxState *GraphicsState) error {
	// This is a comment, so it doesn't change the graphics state or draw anything
	return nil
}

func (ignoreDataBlock *IgnoreDataBlock) String() string {
	return fmt.Sprintf("{COMMENT, %s}", ignoreDataBlock.comment)
}