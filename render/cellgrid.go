package render

import (
	"fmt"
	"image"
	"io"
	"strconv"

	"github.com/otiai10/gat/border"
	"github.com/otiai10/gat/color"
)

// CellGrid ...
type CellGrid struct {

	// Border is a border printer.
	Border border.Border
	// Colorpicker is an algorithm to pickup colors.
	Colorpicker color.Picker

	// Row of cell grid, the terminal height by default if the image is portrait.
	Row uint16
	// Col of cell grid, the terminal width by default if the image is landscape.
	Col uint16
	// Xpixel uint16
	// Ypixel uint16

	// Placeholder is a text printed in the cell, " " by default,
	// ignored in debug mode because color number will be printed instead.
	Placeholder string

	// Debug given, print available color tables
	// and print calculated color number inside the cell.
	// Debug is exclusive property, any other properties might be ignored.
	Debug bool
}

// Render renders specified image by using cell.
func (grid *CellGrid) Render(w io.Writer, img image.Image) error {

	if grid.Row <= 1 || grid.Col <= 1 {
		return fmt.Errorf("output canvas is too small: %+v", grid)
	}

	if grid.Border == nil {
		grid.Border = border.EmptyBorder{}
	}

	rowcount := int(grid.Row - 1)
	// rowcount -= grid.Border.Width()
	for i := 0; i < grid.Border.Width(); i++ {
		rowcount--
	}
	if rowcount <= 0 {
		rowcount = 1
	}
	colcount := int(float64(rowcount) * float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))

	// cell := img.Bounds().Max.Y / rowcount
	cell := float64(img.Bounds().Max.Y) / float64(rowcount)

	// Print top header
	grid.Border.Top(w, colcount+grid.Border.Width())
	if grid.Border.Width() > 0 { // FIXME
		fmt.Fprint(w, "\n")
	}

	for row := 0; row < rowcount; row++ {
		grid.Border.Left(w, row)
		for col := 0; col < colcount; col++ {
			r, g, b, _ := grid.Colorpicker.Pick(img, image.Rectangle{
				Min: image.Point{int(float64(col) * cell), int(float64(row) * cell)},
				Max: image.Point{int(float64(col+1)*cell) - 1, int(float64(row+1)*cell) - 1},
			})
			grid.Fprint(w, color.GetCodeByRGBA(r, g, b, 0))
		}
		grid.Border.Right(w, row)
		fmt.Fprintf(w, "\n")
	}

	// Print bottom footer
	grid.Border.Bottom(w, colcount+grid.Border.Width())

	return nil
}

// Fprint ...
func (grid *CellGrid) Fprint(w io.Writer, code int) {
	if grid.Debug {
		text := "  " + grid.Placeholder + strconv.Itoa(code)
		fmt.Fprintf(w, "\x1b[48;5;%dm%s\x1b[m", code, text[len(text)-3:])
	} else {
		fmt.Fprintf(w, "\x1b[48;5;%dm%s\x1b[m", code, grid.Placeholder)
	}
}
