package render

import "image"

// ITerm renderer to support image printing on iTerm.
// See https://www.iterm2.com/documentation-images.html
type ITerm struct {
}

// Render renders specified image to iTerm stdout.
func (iterm *ITerm) Render(img image.Image) error {
	return nil
}
