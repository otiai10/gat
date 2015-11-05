package main

import (
	"flag"
	"fmt"
	"image"
	"io"
	"log"
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
	log.Println(col, row)
	stdout, stderr := defaultOut, defaultErr
	if len(os.Args) < 2 {
		fmt.Fprint(stderr, "filename required")
		return
	}
	filename := os.Args[1]
	run(filename, stdout, stderr)
}

func colorcheck() {
	for i := 0; i < 256; i++ {
		gcat.Fprint(os.Stdout, i, fmt.Sprintf("%03d", i))
		if i%15 == 0 {
			fmt.Print("\n")
		}
	}

}
func run(filename string, stdout, stderr io.ReadWriter) {
	colorcheck()

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	log.Println(col, row)

	var client *gcat.Client
	switch {
	case col > 0:
		client = gcat.OfTerminal()
	case row > 0:
		client = gcat.OfTerminal()
	default:
		client = gcat.OfTerminal()
	}

	// client := gcat.NewClient()
	client.PrintImage(img)
}
