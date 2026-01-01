package render

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io"
)

// Kitty renderer to support image printing using Kitty Graphics Protocol.
// See https://sw.kovidgoyal.net/kitty/graphics-protocol/
type Kitty struct {
	// Scale size
	Scale float64
}

const (
	// Maximum chunk size for base64 encoded data
	kittyChunkSize = 4096
)

// Render renders specified image using Kitty Graphics Protocol.
func (k *Kitty) Render(w io.Writer, img image.Image) error {
	if k.Scale == 0 {
		k.Scale = 1
	}

	// Encode image as PNG
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		return err
	}

	// Base64 encode the PNG data
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Send in chunks
	// Format: <ESC>_G<control>;<payload><ESC>\
	// Control: a=T (transmit+display), f=100 (PNG), m=1/0 (more/last)
	for len(encoded) > 0 {
		chunk := encoded
		more := 0
		if len(encoded) > kittyChunkSize {
			chunk = encoded[:kittyChunkSize]
			encoded = encoded[kittyChunkSize:]
			more = 1
		} else {
			encoded = ""
		}

		if more == 1 {
			// More chunks to follow
			fmt.Fprintf(w, "\x1b_Ga=T,f=100,m=1;%s\x1b\\", chunk)
		} else {
			// Last chunk (or only chunk)
			fmt.Fprintf(w, "\x1b_Ga=T,f=100,m=0;%s\x1b\\", chunk)
		}
	}

	// Add newline after image
	fmt.Fprintln(w)
	return nil
}

// SetScale sets the scale for the renderer.
func (k *Kitty) SetScale(scale float64) error {
	k.Scale = scale
	return nil
}
