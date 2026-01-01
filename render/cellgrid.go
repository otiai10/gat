package render

import (
	"fmt"
	"image"
	"io"
	"strconv"

	"github.com/otiai10/gat/border"
	"github.com/otiai10/gat/color"
	"github.com/otiai10/gat/terminal"
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

	// HalfBlock enables half-block character mode for 2x vertical resolution.
	// Uses "▀" character with foreground (top) and background (bottom) colors.
	HalfBlock bool

	// TrueColor enables 24-bit RGB colors instead of 256-color palette.
	TrueColor bool
}

// Render renders specified image by using cell.
func (grid *CellGrid) Render(w io.Writer, img image.Image) error {

	if grid.Border == nil {
		grid.Border = border.EmptyBorder{}
	}
	if grid.Colorpicker == nil {
		grid.Colorpicker = color.AverageColorPicker{}
	}
	if grid.Placeholder == "" {
		grid.Placeholder = "  "
	}

	rect := terminal.DefineRectangle(grid.Row, grid.Col, len(grid.Placeholder), img)
	if rect.Row <= 1 || rect.Col <= 1 {
		return fmt.Errorf("output canvas is too small: %+v", rect)
	}

	rowcount := int(rect.Row - 1)
	// rowcount -= grid.Border.Width()
	for i := 0; i < grid.Border.Width(); i++ {
		rowcount--
	}
	if rowcount <= 0 {
		rowcount = 1
	}
	colcount := int(float64(rowcount) * float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))

	// HalfBlock uses single char "▀" instead of placeholder (typically "  "),
	// so we need to scale columns to maintain aspect ratio.
	if grid.HalfBlock {
		colcount *= len(grid.Placeholder)
	}

	// cellHeight is pixel height per terminal row
	cellHeight := float64(img.Bounds().Max.Y) / float64(rowcount)
	// cellWidth is pixel width per terminal column
	cellWidth := float64(img.Bounds().Max.X) / float64(colcount)

	// Print top header
	grid.Border.Top(w, colcount+grid.Border.Width())
	if grid.Border.Width() > 0 { // FIXME
		fmt.Fprint(w, "\n")
	}

	for row := 0; row < rowcount; row++ {
		grid.Border.Left(w, row)
		for col := 0; col < colcount; col++ {
			if grid.HalfBlock {
				// Half-block mode: sample top and bottom halves separately
				halfCellHeight := cellHeight / 2
				topR, topG, topB, _ := grid.Colorpicker.Pick(img, image.Rectangle{
					Min: image.Point{int(float64(col) * cellWidth), int(float64(row)*cellHeight + 0)},
					Max: image.Point{int(float64(col+1)*cellWidth) - 1, int(float64(row)*cellHeight + halfCellHeight)},
				})
				botR, botG, botB, _ := grid.Colorpicker.Pick(img, image.Rectangle{
					Min: image.Point{int(float64(col) * cellWidth), int(float64(row)*cellHeight + halfCellHeight)},
					Max: image.Point{int(float64(col+1)*cellWidth) - 1, int(float64(row+1)*cellHeight) - 1},
				})
				grid.FprintHalfBlock(w, topR, topG, topB, botR, botG, botB)
			} else {
				r, g, b, _ := grid.Colorpicker.Pick(img, image.Rectangle{
					Min: image.Point{int(float64(col) * cellWidth), int(float64(row) * cellHeight)},
					Max: image.Point{int(float64(col+1)*cellWidth) - 1, int(float64(row+1)*cellHeight) - 1},
				})
				if grid.TrueColor {
					grid.FprintTrueColor(w, r, g, b)
				} else {
					grid.Fprint(w, color.GetCodeByRGBA(r, g, b, 0))
				}
			}
		}
		grid.Border.Right(w, row)
		fmt.Fprintf(w, "\n")
	}

	// Print bottom footer
	grid.Border.Bottom(w, colcount+grid.Border.Width())

	return nil
}

// Fprint prints a cell with 256-color palette background.
func (grid *CellGrid) Fprint(w io.Writer, code int) {
	if grid.Debug {
		text := "  " + grid.Placeholder + strconv.Itoa(code)
		fmt.Fprintf(w, "\x1b[48;5;%dm%s\x1b[m", code, text[len(text)-3:])
	} else {
		fmt.Fprintf(w, "\x1b[48;5;%dm%s\x1b[m", code, grid.Placeholder)
	}
}

// FprintTrueColor prints a cell with 24-bit RGB background.
func (grid *CellGrid) FprintTrueColor(w io.Writer, r, g, b uint32) {
	fmt.Fprintf(w, "\x1b[48;2;%d;%d;%dm%s\x1b[m", r>>8, g>>8, b>>8, grid.Placeholder)
}

// FprintHalfBlock prints a half-block character with top (foreground) and bottom (background) colors.
// This effectively renders 2 vertical pixels per terminal cell.
func (grid *CellGrid) FprintHalfBlock(w io.Writer, topR, topG, topB, botR, botG, botB uint32) {
	fmt.Fprintf(w, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀\x1b[m",
		topR>>8, topG>>8, topB>>8,
		botR>>8, botG>>8, botB>>8)
}

// SetScale ...
func (grid *CellGrid) SetScale(scale float64) error {
	return fmt.Errorf("CellGrid renderer doesn't support `scale`")
}
