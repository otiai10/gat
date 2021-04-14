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
	var supportedTerminals = []string{
		"\x1b[?62;", // VT240
		"\x1b[?63;", // wsltty
		"\x1b[?64;", // mintty
		"\x1b[?65;", // RLogin
	}
	supported := false
	for _, supportedTerminal := range supportedTerminals {
		if bytes.HasPrefix(b[:n], []byte(supportedTerminal)) {
			supported = true
			break
		}
	}
	if !supported {
		return false
	}

	sb := b[6:n]
	n = bytes.IndexByte(sb, 'c')
	if n != -1 {
		sb = sb[:n]
	}
	for _, t := range bytes.Split(sb, []byte(";")) {
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
