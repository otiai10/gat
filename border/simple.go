package border

import (
	"fmt"
	"io"
)

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
