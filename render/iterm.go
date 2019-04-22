package render

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"io"
)

// ITerm renderer to support image printing on iTerm.
// See https://www.iterm2.com/documentation-images.html
type ITerm struct {
	// Scale size
	Scale float64
}

// Render renders specified image to iTerm stdout.
func (iterm *ITerm) Render(w io.Writer, img image.Image) error {
	if iterm.Scale == 0 {
		iterm.Scale = 1
	}
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, img)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	size := img.Bounds().Max
	width := int(float64(size.X) * iterm.Scale)
	height := int(float64(size.Y) * iterm.Scale)
	fmt.Fprintf(w, "\033]1337;File=;width=%dpx;height=%dpx;inline=1:%s\a\n", width, height, encoded)
	return nil
}

// SetScale ...
func (iterm *ITerm) SetScale(scale float64) error {
	iterm.Scale = scale
	return nil
}
