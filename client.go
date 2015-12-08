package gat

import (
	"fmt"
	"image"
	"io"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/otiai10/gat/colors"
)

// Cell represents a cell expression of output
var Cell = "  "

// Client ...
type Client struct {
	Out, Err    io.Writer // Out platform
	Canvas      Rect      // output canvas
	Border      Border
	ColorPicker colors.Picker
	IsDebug     bool
}

// NewClient ...
func NewClient(rect Rect) *Client {
	return &Client{
		Out:         os.Stdout,
		Err:         os.Stderr,
		Canvas:      rect,
		Border:      DefaultBorder{},
		ColorPicker: colors.AverageColorPicker{},
	}
}

// Set ...
func (c *Client) Set(attr interface{}) *Client {
	switch attr := attr.(type) {
	case Border:
		c.Border = attr
	case Rect:
		c.Canvas = attr
	case colors.Picker:
		c.ColorPicker = attr
	}
	return c
}

// Debug ...
func (c *Client) Debug(f bool) *Client {
	c.IsDebug = f
	return c
}

// PrintImage ...
func (c *Client) PrintImage(img image.Image) error {
	if c.Canvas.Row <= 1 || c.Canvas.Col <= 1 {
		return fmt.Errorf("output canvas is too small: %+v", c.Canvas)
	}
	rowcount := int(c.Canvas.Row - 1)

	for i := 0; i < c.Border.Width(); i++ {
		rowcount--
	}
	if rowcount <= 0 {
		rowcount = 1
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
			r, g, b, _ := c.ColorPicker.Pick(img, image.Rectangle{
				Min: image.Point{int(float64(col) * cell), int(float64(row) * cell)},
				Max: image.Point{int(float64(col+1)*cell) - 1, int(float64(row+1)*cell) - 1},
			})
			// fmt.Fprintf(c.Out, "%02d", col)
			c.Fprint(c.Out, colors.GetCodeByRGBA(r, g, b, 0), Cell)
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
	// fmt.Fprintf(w, "\x1b[48;5;%dm%s\x1b[m", code, text)
	fmt.Fprintf(w, "\x1b[48;5;%dm%s\x1b[m", code, text)
}

// Fprint ...
func (c *Client) Fprint(w io.Writer, code int, text string) {
	if c.IsDebug {
		text := "  " + Cell + strconv.Itoa(code)
		fmt.Fprintf(w, "\x1b[48;5;%dm%s\x1b[m", code, text[len(text)-3:])
	} else {
		fmt.Fprintf(w, "\x1b[48;5;%dm%s\x1b[m", code, text)
	}
}

// Rect ...
type Rect struct {
	Row uint16
	Col uint16
	// Xpixel uint16
	// Ypixel uint16
}

// GetTerminal ...
func GetTerminal() Rect {
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
	return *t
}
