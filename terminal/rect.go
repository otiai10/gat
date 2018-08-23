package terminal

import (
	"image"
)

// Rect ...
type Rect struct {
	Row uint16
	Col uint16
	// Xpixel uint16
	// Ypixel uint16
}

// DefineRectangle ...
func DefineRectangle(row, col, cellwidth int, img image.Image) Rect {
	switch {
	case row > 0:
		// Define reacangle by given row and the aspect ratio of source image.
		return Rect{
			Row: uint16(row), // restrict output canvas by given "row"
			Col: uint16(float64(row) * (float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))),
		}
	case col > 0:
		// Define reacangle by given col and the aspect ratio of source image.
		return Rect{
			Col: uint16(col),
			Row: uint16(float64(col) * (float64(img.Bounds().Max.Y) / float64(img.Bounds().Max.X))),
		}
	default:
		term := GetTerminal()
		available := (float64(term.Col) / float64(cellwidth)) / float64(term.Row)
		source := float64(img.Bounds().Size().X) / float64(img.Bounds().Size().Y)
		if source > available {
			term.Row = uint16(float64(term.Col) / source / float64(cellwidth))
		} else {
			term.Col = uint16(float64(term.Row) * source / float64(cellwidth))
		}
		return term
	}
}
