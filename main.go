package main

import (
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"syscall"
	"unsafe"

	_ "image/png"
)

var (
	defaultOut = os.Stdout
	defaultErr = os.Stderr
)

func main() {
	run(defaultOut, defaultErr)
}

func run(stdout, stderr io.ReadWriter) {
	t := getTerminal()
	for r := 0; r < int(t.Row); r++ {
		for c := 0; c < int(t.Col)-1; c++ {
			fmt.Printf("\x1b[48;5;%dm \x1b[m", r)
			// fmt.Print("%K[42m%k") not work
			// fmt.Print("\x1b%K{42}%k")
		}
		fmt.Print("#\n")
	}

	f, err := os.Open("gopher.png")
	if err != nil {
		panic(err)
	}

	img, format, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	log.Println(img.Bounds().Max.X, img.Bounds().Max.Y, format)
	log.Println(t.Col, t.Xpixel, t.Row, t.Ypixel)

	ratio := 4

	for row := 0; row < img.Bounds().Max.Y/ratio; row++ {
		for col := 0; col < img.Bounds().Max.X/ratio; col++ {
			r, g, b, _ := img.At(col*ratio, row*ratio).RGBA()
			if r == 0 && g == 0 && b == 0 {
				fmt.Printf("\x1b[48;5;%dm \x1b[m", 0)
			} else {
				fmt.Printf("\x1b[48;5;%[1]dm \x1b[m", r)
			}
		}
		fmt.Print("\n")
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
