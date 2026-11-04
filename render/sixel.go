package render

import (
	"image"
	"io"

	sixel "github.com/mattn/go-sixel"
	"golang.org/x/image/draw"
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
	// Like the previous nfnt/resize.Thumbnail implementation, only shrink, never enlarge.
	if sxl.Scale < 1 {
		size := img.Bounds().Size()
		width := max(1, int(float64(size.X)*sxl.Scale))
		height := max(1, int(float64(size.Y)*sxl.Scale))
		dst := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Src, nil)
		img = dst
	}
	return sixel.NewEncoder(w).Encode(img)
}

// SetScale ...
func (sxl *Sixel) SetScale(scale float64) error {
	sxl.Scale = scale
	return nil
}
