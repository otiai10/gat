package gcat

import (
	"fmt"
	"image"
	"io"
	"os"
	"syscall"
	"unsafe"

	"github.com/otiai10/gcat/colors"
)

// Client ...
type Client struct {
	Output, Errput io.ReadWriter // Output platform
	Canvas         Rect          // output canvas
}

func NewClient(rect Rect) *Client {
	return &Client{
		os.Stdout,
		os.Stderr,
		rect,
	}
}

func OfTerminal() *Client {
	t := GetTerminal()
	return NewClient(Rect{
		Row: t.Row,
		Col: t.Col,
	})
}

func (client *Client) PrintImage(img image.Image) error {
	rowcount := int(client.Canvas.Row)
	colcount := int(float64(rowcount) * float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))
	// ratio := img.Bounds().Max.Y / rowcount
	cell := img.Bounds().Max.Y / rowcount
	aspect := 2
	// for row := 0; row < img.Bounds().Max.Y/ratio; row++ {
	for row := 0; row < rowcount; row++ {
		// for col := 0; col < img.Bounds().Max.X/ratio; col++ {
		for col := 0; col < colcount; col++ {
			r, g, b, a := img.At(col*cell+2, row*cell+2).RGBA() // FIXME: 微調整
			for i := 0; i < aspect; i++ {
				// fmt.Fprintf(stdout, "\x1b[48;5;%sm \x1b[m", colors.GetCodeByRGBA(r, g, b, a))
				Fprint(client.Output, colors.GetCodeByRGBA(r, g, b, a), " ")
			}
		}
		fmt.Fprint(client.Output, "\n")
	}
	return nil
}

func (client *Client) Print() error {
	return nil
}

func Fprint(w io.Writer, code int, text string) {
	fmt.Fprintf(w, "\x1b[48;5;%dm%s\x1b[m", code, text)
}

type Rect struct {
	Row uint16
	Col uint16
	// Xpixel uint16
	// Ypixel uint16
}

func GetTerminal() *Rect {
	t := new(Rect)
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
