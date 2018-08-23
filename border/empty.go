package border

import "io"

// EmptyBorder display nothing for borders.
type EmptyBorder struct{}

// Top do nothing.
func (emptyborder EmptyBorder) Top(out io.Writer, cols int) {}

// Left do nothing.
func (emptyborder EmptyBorder) Left(out io.Writer, row int) {}

// Right do nothing.
func (emptyborder EmptyBorder) Right(out io.Writer, row int) {}

// Bottom do nothing.
func (emptyborder EmptyBorder) Bottom(out io.Writer, cols int) {}

// Width claim no width.
func (emptyborder EmptyBorder) Width() int { return 0 }
