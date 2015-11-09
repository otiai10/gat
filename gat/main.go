package main

import (
	"flag"
	"image"
	"io"
	"os"

	"github.com/otiai10/gat"

	_ "image/png"
)

var (
	defaultOut    = os.Stdout
	defaultErr    = os.Stderr
	border, debug bool
)

func init() {
}

func main() {
	/*
		var w = flag.Int("w", 0, "cols")
		var h = flag.Int("h", 0, "rows")
	*/
	flag.BoolVar(&border, "b", false, "border style")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()
	stdout, stderr := defaultOut, defaultErr
	filename := flag.Arg(0)
	run(filename, stdout, stderr, 0, 0)
}

func run(filename string, stdout, stderr io.ReadWriter, col, row int) {

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	var client *gat.Client
	switch {
	case col > 0:
		client = gat.NewClient(gat.Rect{
			Col: uint16(col),
			Row: uint16(col * (img.Bounds().Max.Y / img.Bounds().Max.X)),
		})
	case row > 0:
		client = gat.NewClient(gat.Rect{
			Row: uint16(row),
			Col: uint16(row * (img.Bounds().Max.X / img.Bounds().Max.Y)),
		})
	default:
		client = gat.Terminal()
	}

	client.Out = stdout
	client.Err = stderr

	switch {
	case debug:
		client.Set(gat.DebugBorder{})
	case border:
		client.Set(gat.SimpleBorder{})
	}

	// client := gat.NewClient()
	client.PrintImage(img)
}
