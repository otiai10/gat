package main

import (
	"flag"
	"image"
	"io"
	"log"
	"os"

	"github.com/otiai10/gat/border"

	"github.com/otiai10/gat/color"
	"github.com/otiai10/gat/render"
	"github.com/otiai10/gat/terminal"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	debug       = false
	h           = 0
	w           = 0
	printborder = false
	placeholder = "  "
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Print debug information")
	flag.IntVar(&h, "H", 0, "Height of output")
	flag.IntVar(&w, "W", 0, "Width of output")
	flag.BoolVar(&printborder, "b", false, "Print border")
	flag.StringVar(&placeholder, "s", "  ", "Placeholder text for grid cell")
	flag.Parse()
}

func main() {
	stdout, stderr := os.Stdout, os.Stderr
	filename := flag.Arg(0)
	run(filename, stdout, stderr, h, w)
}

func run(filename string, stdout, stderr io.Writer, row, col int) {

	if debug {
		color.Dump(stderr)
		if filename == "" {
			return
		}
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalln(err)
	}

	rect := defineRectangle(row, col, len(placeholder), img)

	renderer := &render.CellGrid{
		Row:         rect.Row,
		Col:         rect.Col,
		Colorpicker: color.AverageColorPicker{},
		Placeholder: placeholder,
		Debug:       debug,
	}
	if printborder {
		renderer.Border = border.SimpleBorder{}
	}

	if err := renderer.Render(stdout, img); err != nil {
		log.Fatalln(err)
	}
}

func defineRectangle(row, col, cellwidth int, img image.Image) terminal.Rect {
	switch {
	case row > 0:
		// Define reacangle by given row and the aspect ratio of source image.
		return terminal.Rect{
			Row: uint16(row), // restrict output canvas by given "row"
			Col: uint16(float64(row) * (float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))),
		}
	case col > 0:
		// Define reacangle by given col and the aspect ratio of source image.
		return terminal.Rect{
			Col: uint16(col),
			Row: uint16(float64(col) * (float64(img.Bounds().Max.Y) / float64(img.Bounds().Max.X))),
		}
	default:
		term := terminal.GetTerminal()
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
