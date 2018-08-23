package render

import (
	"image"
	"io"
	"os"

	sixel "github.com/mattn/go-sixel"
	"github.com/nfnt/resize"
)

// Sixel ...
// See https://github.com/saitoha/libsixel
type Sixel struct {
	Scale float64
}

// Render renders specified image to iTerm stdout.
func (sxl *Sixel) Render(w io.Writer, img image.Image) error {
	size := img.Bounds().Max
	width := uint(float64(size.X) * sxl.Scale)
	height := uint(float64(size.Y) * sxl.Scale)
	resized := resize.Thumbnail(width, height, img, resize.Bicubic)
	return sixel.NewEncoder(os.Stdout).Encode(resized)
}
