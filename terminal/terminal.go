package terminal

import (
	"os"

	"golang.org/x/term"
)

const (
	// fallbackCols and fallbackRows are used when the terminal size cannot be determined,
	// e.g. when neither stdout nor stdin is attached to a terminal.
	fallbackCols = 80
	fallbackRows = 24
)

// GetTerminal returns the size of the current terminal.
// It queries stdout first, then stdin, and falls back to 80x24
// if neither of them is a terminal.
func GetTerminal() Rect {
	for _, f := range []*os.File{os.Stdout, os.Stdin} {
		if cols, rows, err := term.GetSize(int(f.Fd())); err == nil && cols > 0 && rows > 0 {
			return Rect{Row: uint16(rows), Col: uint16(cols)}
		}
	}
	return Rect{Row: fallbackRows, Col: fallbackCols}
}
