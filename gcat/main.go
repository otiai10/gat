package main

import (
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
)

func main() {
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

	// client := gcat.NewClient()
	gcat.OfTerminal().PrintImage(img)

	// TODO: zoom
	// TODO: or auto zoom by terminal
	/*
		aspect := 2
		t := gcat.GetTerminal()
		rowcount := int(t.Row)
		// colcount := rowcount * (img.Bounds().Max.X / img.Bounds().Max.Y)
		colcount := int(float64(rowcount) * float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))
		ratio := img.Bounds().Max.Y / rowcount
	*/
	/*
		// for row := 0; row < img.Bounds().Max.Y/ratio; row++ {
		for row := 0; row < rowcount; row++ {
			// for col := 0; col < img.Bounds().Max.X/ratio; col++ {
			for col := 0; col < colcount; col++ {
				r, g, b, a := img.At(col*ratio+2, row*ratio+2).RGBA() // FIXME: 微調整
				for i := 0; i < aspect; i++ {
					// fmt.Fprintf(stdout, "\x1b[48;5;%sm \x1b[m", colors.GetCodeByRGBA(r, g, b, a))
					gcat.Fprint(stdout, colors.GetCodeByRGBA(r, g, b, a), " ")
				}
			}
			fmt.Fprint(stdout, "\n")
		}
	*/
}
