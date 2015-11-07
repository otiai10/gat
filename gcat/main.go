package main

import (
	"flag"
	"image"
	"io"
	"os"

	"github.com/otiai10/gcat"

	_ "image/png"
)

var (
	defaultOut = os.Stdout
	defaultErr = os.Stderr
)

func init() {
}

func main() {
	/*
		var w = flag.Int("w", 0, "cols")
		var h = flag.Int("h", 0, "rows")
	*/
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

	var client *gcat.Client
	switch {
	/*
		case col > 0:
			client = gcat.Terminal()
		case row > 0:
			client = gcat.Terminal()
	*/
	default:
		client = gcat.Terminal()
	}

	// client := gcat.NewClient()
	client.PrintImage(img)
}
