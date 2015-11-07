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
	Border         Border
}

// NewClient ...
func NewClient(rect Rect) *Client {
	return &Client{
		Output: os.Stdout,
		Errput: os.Stderr,
		Canvas: rect,
		Border: DefaultBorder{},
		// Border: DebugBorder{},
		// Border: SimpleBorder{},
	}
}

// Terminal ...
func Terminal() *Client {
	t := getTerminal()
	return NewClient(Rect{
		Row: t.Row,
		Col: t.Col,
	})
}

// Set ...
func (client *Client) Set(attr interface{}) *Client {
	switch attr := attr.(type) {
	case Border:
		client.Border = attr
	case Rect:
		client.Canvas = attr
	}
	return client
}

// PrintImage ...
func (client *Client) PrintImage(img image.Image) error {
	rowcount := int(client.Canvas.Row - 1)

	for i := 0; i < client.Border.Width(); i++ {
		rowcount--
	}

	colcount := int(float64(rowcount) * float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))

	// ratio := img.Bounds().Max.Y / rowcount
	cell := img.Bounds().Max.Y / rowcount

	// Print top header
	for col := 0; col < colcount+client.Border.Width(); col++ {
		client.Border.Top(client.Output, col)
	}
	if client.Border.Width() > 0 { // FIXME
		fmt.Fprint(client.Output, "\n")
	}

	for row := 0; row < rowcount; row++ {
		client.Border.Left(client.Output, row)
		for col := 0; col < colcount; col++ {
			r, g, b, a := img.At(col*cell+2, row*cell+2).RGBA() // FIXME: 微調整
			// fmt.Fprintf(client.Output, "%02d", col)
			Fprint(client.Output, colors.GetCodeByRGBA(r, g, b, a), "  ")
		}
		client.Border.Right(client.Output, row)
		// fmt.Fprintf(client.Output, "\n%02d", row)
		fmt.Fprintf(client.Output, "\n")
	}

	// Print bottom footer
	for col := 0; col < colcount+client.Border.Width(); col++ {
		client.Border.Bottom(client.Output, col)
	}

	return nil
}

// Fprint ...
func Fprint(w io.Writer, code int, text string) {
	fmt.Fprintf(w, "\x1b[48;5;%dm%s\x1b[m", code, text)
}

// Rect ...
type Rect struct {
	Row uint16
	Col uint16
	// Xpixel uint16
	// Ypixel uint16
}

func getTerminal() *Rect {
	t := new(Rect)
	retCode, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(t)),
	)

	if int(retCode) == -1 {
		panic(err)
	}
	return t
}
