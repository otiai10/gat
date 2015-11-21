package main

import (
	"flag"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/otiai10/gat"
	"github.com/otiai10/gat/colors"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	defaultOut    = os.Stdout
	defaultErr    = os.Stderr
	border, debug bool
	w, h          int
	picker, cell  string
)

func init() {
	flag.IntVar(&w, "w", 0, descriptionWidth)
	flag.IntVar(&h, "h", 0, descriptionHeight)
	flag.BoolVar(&border, "b", false, descriptionBorder)
	flag.StringVar(&picker, "picker", "average", descriptionPicker)
	flag.StringVar(&cell, "cell", "  ", descriptionCell)
	flag.BoolVar(&debug, "debug", false, descriptionDebug)
	flag.Parse()
}

func main() {
	stdout, stderr := defaultOut, defaultErr
	filename := flag.Arg(0)
	run(filename, stdout, stderr, w, h)
}

func onerror(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}

func run(filename string, stdout, stderr io.ReadWriter, col, row int) {

	f, err := os.Open(filename)
	onerror(err)

	img, _, err := image.Decode(f)
	onerror(err)

	gat.Cell = cell
	client := &gat.Client{}

	switch {
	case debug:
		gat.Cell = "   " // with length 3, to print 3 digit color code in cell.
		colors.Check(stdout)
		client.Set(gat.DebugBorder{Padding: gat.Cell}).Debug(true)
	case border:
		client.Set(gat.SimpleBorder{})
	default:
		client.Set(gat.DefaultBorder{})
	}

	switch picker {
	case "center":
		client.Set(colors.CenterColorPicker{})
	case "lefttop":
		client.Set(colors.LeftTopColorPicker{})
	case "horizontal":
		client.Set(colors.HorizontalAverageColorPicker{})
	default:
		client.Set(colors.AverageColorPicker{})
	}

	client.Out = stdout
	client.Err = stderr

	client.Canvas = getCanvas(col, row, img)

	onerror(client.PrintImage(img))
}

func getCanvas(col, row int, img image.Image) gat.Rect {
	switch {
	case col > 0:
		return gat.Rect{
			Col: uint16(col), // restrict output canvas by given "col"
			Row: uint16(float64(col) * (float64(img.Bounds().Max.Y) / float64(img.Bounds().Max.X))),
		}
	case row > 0:
		return gat.Rect{
			Row: uint16(row), // restrict output canvas by given "row"
			Col: uint16(float64(row) * (float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))),
		}
	default:
		canvas, terminal := gat.Rect{}, gat.GetTerminal()
		rAvailable, rSource := float64(terminal.Col)/float64(terminal.Row)/float64(len(gat.Cell)), float64(img.Bounds().Size().X)/float64(img.Bounds().Size().Y)
		if rAvailable > rSource { // source image is vertically bigger than available canvas
			canvas.Row = terminal.Row // restrict output canvas by current terminal's row
			canvas.Col = uint16(float64(terminal.Row) * rSource / float64(len(gat.Cell)))
		} else { // source image is horizontally bigger than available canvas
			canvas.Col = terminal.Col // restrict output canvas by current terminal's col
			canvas.Row = uint16(float64(terminal.Col) / rSource / float64(len(gat.Cell)))
		}
		return canvas
	}
}

const (
	descriptionCell   = `Cell for output. Output would be constructed with this text.`
	descriptionWidth  = `Width of output canvas. Current terminal width in default.`
	descriptionHeight = `Height of output canvas. Current terminal height in default.`
	descriptionBorder = `Set border to output, such as ╔════════════════════╗`
	descriptionPicker = `Set color picker. ["average" | "horizontal" | "center" | "lefttop"]`
	descriptionDebug  = `Show debug output and debug border.`
)
