package border

import (
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
