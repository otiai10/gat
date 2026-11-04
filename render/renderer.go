package render

import (
	"bytes"
	"image"
	"io"
	"os"
	"time"

	"golang.org/x/term"
)

// Renderer interface.
type Renderer interface {
	Render(io.Writer, image.Image) error
	SetScale(float64) error
}

// ITermImageSupported reports whether the terminal is iTerm2 and can show inline images.
// Inside tmux, TERM_PROGRAM is overwritten by tmux, so LC_TERMINAL (set by iTerm2) is also checked,
// and tmux must allow passthrough so that the image sequence reaches iTerm2.
func ITermImageSupported() bool {
	iterm := os.Getenv("TERM_PROGRAM") == "iTerm.app"
	if !inTmux() {
		return iterm
	}
	return (iterm || os.Getenv("LC_TERMINAL") == "iTerm2") && tmuxPassthroughAllowed()
}

// KittySupported checks if the terminal supports Kitty Graphics Protocol.
// Supported terminals: Kitty, WezTerm, Ghostty
func KittySupported() bool {
	// Check for Kitty terminal
	if os.Getenv("KITTY_WINDOW_ID") != "" {
		return true
	}
	// Check for Ghostty
	if os.Getenv("TERM_PROGRAM") == "ghostty" {
		return true
	}
	// Check for WezTerm (also supports Kitty protocol)
	if os.Getenv("TERM_PROGRAM") == "WezTerm" {
		return true
	}
	return false
}

// sixelQueryTimeout is how long SixelSupported waits for the terminal to answer.
const sixelQueryTimeout = 500 * time.Millisecond

// SixelSupported reports whether the terminal attached to stdout supports Sixel graphics,
// by sending a DA1 (Primary Device Attributes) query.
func SixelSupported() bool {
	fd := int(os.Stdout.Fd())
	if !term.IsTerminal(fd) {
		return false
	}
	s, err := term.MakeRaw(fd)
	if err != nil {
		return false
	}
	defer term.Restore(fd, s)
	if _, err := os.Stdout.Write([]byte("\x1b[c")); err != nil {
		return false
	}
	// Do not block forever if the terminal does not answer the DA1 query.
	// SetReadDeadline may be unsupported on some platforms; then it just blocks as before.
	_ = os.Stdout.SetReadDeadline(time.Now().Add(sixelQueryTimeout))
	defer os.Stdout.SetReadDeadline(time.Time{})

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
	for t := range bytes.SplitSeq(sb, []byte(";")) {
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
	case KittySupported():
		return &Kitty{Scale: 1}
	case SixelSupported():
		return &Sixel{Scale: 1}
	default:
		// Same cell-grid settings as the gat CLI: half-block characters in 24-bit color.
		return &CellGrid{HalfBlock: true, TrueColor: true}
	}
}
