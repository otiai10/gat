package gat

import (
	"fmt"
	"io"
)

// Border defines border decoration.
type Border interface {
	Top(out io.Writer, cols int)
	Left(out io.Writer, row int)
	Right(out io.Writer, row int)
	Bottom(out io.Writer, cols int)
	Width() int
}

// DefaultBorder display nothing for borders.
type DefaultBorder struct{}

// Top do nothing.
func (border DefaultBorder) Top(out io.Writer, cols int) {}

// Left do nothing.
func (border DefaultBorder) Left(out io.Writer, row int) {}

// Right do nothing.
func (border DefaultBorder) Right(out io.Writer, row int) {}

// Bottom do nothing.
func (border DefaultBorder) Bottom(out io.Writer, cols int) {}

// Width claim no width.
func (border DefaultBorder) Width() int {
	return 0
}

// DebugBorder ...
type DebugBorder struct {
	Padding string
}

// Top ...
func (border DebugBorder) Top(out io.Writer, cols int) {
	fmt.Fprint(out, "  ") // keep space to row index on the left edge
	for c := 1; c < cols-1; c++ {
		s := "  " + border.Padding + fmt.Sprintf("%d", c-1)
		fmt.Fprintf(out, s[len(s)-3:])
	}
}

// Left ...
func (border DebugBorder) Left(out io.Writer, row int) {
	s := "  " + fmt.Sprintf("%d", row)
	fmt.Fprintf(out, s[len(s)-2:])
}

// Right ...
func (border DebugBorder) Right(out io.Writer, row int) {
}

// Bottom ...
func (border DebugBorder) Bottom(out io.Writer, cols int) {
}

// Width ...
func (border DebugBorder) Width() int {
	return 2
}

// SimpleBorder ...
type SimpleBorder struct{}

// Top ...
func (border SimpleBorder) Top(out io.Writer, cols int) {
	fmt.Fprintf(out, "╔═")
	for c := 1; c < cols-1; c++ {
		fmt.Fprintf(out, "══")
	}
	fmt.Fprintf(out, "═╗")
}

// Left ...
func (border SimpleBorder) Left(out io.Writer, row int) {
	fmt.Fprintf(out, "║ ")
}

// Right ...
func (border SimpleBorder) Right(out io.Writer, row int) {
	fmt.Fprintf(out, " ║")
}

// Bottom ...
func (border SimpleBorder) Bottom(out io.Writer, cols int) {
	fmt.Fprintf(out, "╚═")
	for c := 1; c < cols-1; c++ {
		fmt.Fprintf(out, "══")
	}
	fmt.Fprintf(out, "═╝")
}

// Width ...
func (border SimpleBorder) Width() int {
	return 2
}
