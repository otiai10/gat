package main

import (
	"flag"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/otiai10/gcat"

	_ "image/png"
)

var (
	defaultOut = os.Stdout
	defaultErr = os.Stderr
	col, row   int
)

func init() {
}

func main() {
	flag.IntVar(&col, "col", 0, "col")
	flag.IntVar(&row, "row", 0, "row")
	flag.Parse()
	stdout, stderr := defaultOut, defaultErr
	if len(os.Args) < 2 {
		fmt.Fprint(stderr, "filename required")
		return
	}
	filename := os.Args[1]
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
	case col > 0:
		client = gcat.Terminal()
	case row > 0:
		client = gcat.Terminal()
	default:
		client = gcat.Terminal()
	}

	// client := gcat.NewClient()
	client.PrintImage(img)
}
