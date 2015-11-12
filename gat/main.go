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
			Col: uint16(col),
			Row: uint16(float64(col) * (float64(img.Bounds().Max.Y) / float64(img.Bounds().Max.X))),
		})
	case row > 0:
		client = gat.NewClient(gat.Rect{
			Row: uint16(row),
			Col: uint16(float64(row) * (float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))),
		})
	default:
		client = gat.Terminal()
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

	// client := gat.NewClient()
	client.PrintImage(img)
}
