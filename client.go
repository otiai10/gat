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
	Out, Err    io.ReadWriter // Out platform
	Canvas      Rect          // output canvas
	Border      Border
	ColorPicker colors.Picker
}

// NewClient ...
func NewClient(rect Rect) *Client {
	return &Client{
		Out:         os.Stdout,
		Err:         os.Stderr,
		Canvas:      rect,
		Border:      DefaultBorder{},
		ColorPicker: colors.AverageColor,
		// ColorPicker: colors.AverageColorX,
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

	// cell := img.Bounds().Max.Y / rowcount
	cell := float64(img.Bounds().Max.Y) / float64(rowcount)

	// Print top header
	c.Border.Top(c.Out, colcount+c.Border.Width())
	if c.Border.Width() > 0 { // FIXME
		fmt.Fprint(c.Out, "\n")
	}

	for row := 0; row < rowcount; row++ {
		c.Border.Left(c.Out, row)
		for col := 0; col < colcount; col++ {
			r, g, b, _ := c.ColorPicker(img, image.Rectangle{
				Min: image.Point{int(float64(col) * cell), int(float64(row) * cell)},
				Max: image.Point{int(float64(col+1)*cell) - 1, int(float64(row+1)*cell) - 1},
			})
			// fmt.Fprintf(c.Out, "%02d", col)
			Fprint(c.Out, colors.GetCodeByRGBA(r, g, b, 0), "  ")
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
