package main

import (
	"flag"
	"image"
	"io"
	"log"
	"os"

	"github.com/otiai10/gat/color"
	"github.com/otiai10/gat/render"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	debug               = false
	h                   = 0
	w                   = 0
	printborder         = false
	placeholder         = "  "
	scale       float64 = 1
	usecell             = false
)

func init() {
	flag.BoolVar(&debug, "debug", false, "Print debug information")
	flag.IntVar(&h, "H", 0, "Height of output")
	flag.IntVar(&w, "W", 0, "Width of output")
	flag.Float64Var(&scale, "S", 1, "Scale for iTerm image output")
	flag.BoolVar(&printborder, "b", false, "Print border")
	flag.StringVar(&placeholder, "t", "  ", "Placeholder text for grid cell")
	flag.BoolVar(&usecell, "c", false, "Prefer cell grid output than terminal app")
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

	renderer := getRenderer(usecell, row, col, placeholder, scale, img)

	if err := renderer.Render(stdout, img); err != nil {
		log.Fatalln(err)
	}
}

func getRenderer(usecell bool, row, col int, placeholder string, scale float64, img image.Image) render.Renderer {
	switch {
	case !usecell && render.ITermImageSupported():
		return &render.ITerm{
			Scale: scale,
		}
	case !usecell && render.SixelSupported():
		return &render.Sixel{
			Scale: scale,
		}
	default:
		return &render.CellGrid{
			Row:         uint16(row),
			Col:         uint16(col),
			Colorpicker: color.AverageColorPicker{},
			Placeholder: placeholder,
		}
	}
}
