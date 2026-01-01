package border

import (
	"fmt"
	"io"
)

// DebugBorder ...
type DebugBorder struct {
	Padding string
}

// Top ...
func (border DebugBorder) Top(out io.Writer, cols int) {
	fmt.Fprint(out, "  ") // keep space to row index on the left edge
	for c := 1; c < cols-1; c++ {
		s := "  " + border.Padding + fmt.Sprintf("%d", c-1)
		fmt.Fprint(out, s[len(s)-3:])
	}
}

// Left ...
func (border DebugBorder) Left(out io.Writer, row int) {
	s := "  " + fmt.Sprintf("%d", row)
	fmt.Fprint(out, s[len(s)-2:])
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
