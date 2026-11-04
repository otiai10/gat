package render

import (
	"image"
	"io"

	sixel "github.com/mattn/go-sixel"
)

// Sixel renderer to support image printing on Sixel-capable terminals.
// See https://github.com/saitoha/libsixel
type Sixel struct {
	Scale float64
}

// Render renders specified image to w as Sixel graphics.
func (sxl *Sixel) Render(w io.Writer, img image.Image) error {
	if sxl.Scale == 0 {
		sxl.Scale = 1
	}
	img = shrink(img, sxl.Scale)
	return sixel.NewEncoder(w).Encode(img)
}

// SetScale ...
func (sxl *Sixel) SetScale(scale float64) error {
	sxl.Scale = scale
	return nil
}
