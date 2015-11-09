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
	run(filename, stdout, stderr)
}

func run(filename string, stdout, stderr io.ReadWriter) {

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
	/*
		case col > 0:
			client = gat.Terminal()
		case row > 0:
			client = gat.Terminal()
	*/
	default:
		client = gat.Terminal()
	}

	switch {
	case debug:
		client.Set(gat.DebugBorder{})
	case border:
		client.Set(gat.SimpleBorder{})
	}

	// client := gat.NewClient()
	client.PrintImage(img)
}
