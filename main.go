package main

import (
	"fmt"
	"image"
	"io"
	"os"
	"syscall"
	"unsafe"

	"github.com/otiai10/gcat/colors"

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

func run(filename string, stdout, stderr io.ReadWriter) {
	for i := 0; i < 256; i++ {
		fmt.Printf("\x1b[48;5;%dm%03d\x1b[m", i, i)
		if i%15 == 0 {
			fmt.Print("\n")
		}
	}

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	// TODO: zoom
	// TODO: or auto zoom by terminal
	aspect := 2
	t := getTerminal()
	rowcount := int(t.Row)
	// colcount := rowcount * (img.Bounds().Max.X / img.Bounds().Max.Y)
	colcount := int(float64(rowcount) * float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))
	ratio := img.Bounds().Max.Y / rowcount

	// for row := 0; row < img.Bounds().Max.Y/ratio; row++ {
	for row := 0; row < rowcount; row++ {
		// for col := 0; col < img.Bounds().Max.X/ratio; col++ {
		for col := 0; col < colcount; col++ {
			r, g, b, a := img.At(col*ratio+2, row*ratio+2).RGBA() // FIXME: 微調整
			for i := 0; i < aspect; i++ {
				fmt.Fprintf(stdout, "\x1b[48;5;%sm \x1b[m", colors.GetCodeByRGBA(r, g, b, a))
			}
		}
		fmt.Fprint(stdout, "\n")
	}
}

type terminal struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getTerminal() *terminal {
	t := new(terminal)
	retCode, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(t)),
	)

	if int(retCode) == -1 {
		panic(errno)
	}
	return t
}
