package gat

import (
	"fmt"
	"image"
	"io"
	"os"
	"syscall"
	"unsafe"

	"github.com/otiai10/gat/colors"
)

// Client ...
type Client struct {
	Out, Err io.ReadWriter // Out platform
	Canvas   Rect          // output canvas
	Border   Border
}

// NewClient ...
func NewClient(rect Rect) *Client {
	return &Client{
		Out:    os.Stdout,
		Err:    os.Stderr,
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
func (c *Client) Set(attr interface{}) *Client {
	switch attr := attr.(type) {
	case Border:
		c.Border = attr
	case Rect:
		c.Canvas = attr
	}
	return c
}

// PrintImage ...
func (c *Client) PrintImage(img image.Image) error {
	rowcount := int(c.Canvas.Row - 1)

	for i := 0; i < c.Border.Width(); i++ {
		rowcount--
	}

	colcount := int(float64(rowcount) * float64(img.Bounds().Max.X) / float64(img.Bounds().Max.Y))

	// ratio := img.Bounds().Max.Y / rowcount
	cell := img.Bounds().Max.Y / rowcount

	// Print top header
	c.Border.Top(c.Out, colcount+c.Border.Width())
	if c.Border.Width() > 0 { // FIXME
		fmt.Fprint(c.Out, "\n")
	}

	for row := 0; row < rowcount; row++ {
		c.Border.Left(c.Out, row)
		for col := 0; col < colcount; col++ {
			r, g, b, a := img.At(col*cell+2, row*cell+2).RGBA() // FIXME: 微調整
			// fmt.Fprintf(c.Out, "%02d", col)
			Fprint(c.Out, colors.GetCodeByRGBA(r, g, b, a), "  ")
		}
		c.Border.Right(c.Out, row)
		// fmt.Fprintf(c.Out, "\n%02d", row)
		fmt.Fprintf(c.Out, "\n")
	}

	// Print bottom footer
	c.Border.Bottom(c.Out, colcount+c.Border.Width())

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
