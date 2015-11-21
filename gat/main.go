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
)

func init() {
	flag.IntVar(&w, "w", 0, "Width of output (cols)")
	flag.IntVar(&h, "h", 0, "Height of output (rows)")
	flag.BoolVar(&border, "b", false, "Add border")
	flag.BoolVar(&debug, "debug", false, "Debug mode with cell index")
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

	var client *gat.Client
	switch {
	case col > 0:
		client = gat.NewClient(gat.Rect{
			Col: uint16(col), // restrict output canvas by given "col"
			Row: uint16(float64(col) * (float64(img.Bounds().Max.Y) / float64(img.Bounds().Max.X))),
		})
	case row > 0:
		client = gat.NewClient(gat.Rect{
			Row: uint16(row), // restrict output canvas by given "row"
			Col: uint16(float64(row) * (float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))),
		})
	default:
		canvas, terminal := gat.Rect{}, gat.GetTerminal()
		rAvailable, rSource := float64(terminal.Col)/float64(terminal.Row), float64(img.Bounds().Size().X)/float64(img.Bounds().Size().Y)
		if rAvailable > rSource { // source image is vertically bigger than available canvas
			canvas.Row = terminal.Row // restrict output canvas by current terminal's row
			canvas.Col = uint16(float64(terminal.Row) * rSource / float64(len(gat.Cell)))
		} else { // source image is horizontally bigger than available canvas
			canvas.Col = terminal.Col // restrict output canvas by current terminal's col
			canvas.Row = uint16(float64(terminal.Col) / rSource / float64(len(gat.Cell)))
		}
		client = gat.NewClient(canvas)
	}

	client.Out = stdout
	client.Err = stderr

	switch {
	case debug:
		colors.Check(stdout)
		client.Set(gat.DebugBorder{})
	case border:
		client.Set(gat.SimpleBorder{})
	}

	onerror(client.PrintImage(img))
}
