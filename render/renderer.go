package render

import (
	"bytes"
	"image"
	"io"
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

// Renderer interface.
type Renderer interface {
	Render(io.Writer, image.Image) error
	SetScale(float64) error
}

// ITermImageSupported ...
func ITermImageSupported() bool {
	return os.Getenv("TERM_PROGRAM") == "iTerm.app"
}

// SixelSupported ...
func SixelSupported() bool {
	s, err := terminal.MakeRaw(1)
	if err != nil {
		return false
	}
	defer terminal.Restore(1, s)
	_, err = os.Stdout.Write([]byte("\x1b[c"))
	if err != nil {
		return false
	}
	defer readTimeout(os.Stdout, time.Time{})

	var b [100]byte
	n, err := os.Stdout.Read(b[:])
	if err != nil {
		return false
	}
	if !bytes.HasPrefix(b[:n], []byte("\x1b[?63;")) {
		return false
	}
	for _, t := range bytes.Split(b[4:n], []byte(";")) {
		if len(t) == 1 && t[0] == '4' {
			return true
		}
	}
	return false
}

// GetDefaultRenderer provides an applicable renderer for current platform.
func GetDefaultRenderer() Renderer {
	switch {
	case ITermImageSupported():
		return &ITerm{Scale: 1}
	case SixelSupported():
		return &Sixel{Scale: 1}
	default:
		return &CellGrid{}
	}
}
