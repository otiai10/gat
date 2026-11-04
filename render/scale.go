package render

import (
	"image"

	"golang.org/x/image/draw"
)

// shrink returns img scaled down by scale. Like the previous
// nfnt/resize.Thumbnail implementation, it only shrinks, never enlarges:
// for scale >= 1 the original image is returned as is.
func shrink(img image.Image, scale float64) image.Image {
	if scale >= 1 {
		return img
	}
	size := img.Bounds().Size()
	width := max(1, int(float64(size.X)*scale))
	height := max(1, int(float64(size.Y)*scale))
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Src, nil)
	return dst
}
