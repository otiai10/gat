package colors

import "image"

// Picker ...
type Picker interface {
	Pick(image.Image, image.Rectangle) (r, g, b, a uint32)
}

// AverageColorPicker picks average color of given RectAngle area of src image.
type AverageColorPicker struct{}

// Pick of AverageColorPicker.
func (picker AverageColorPicker) Pick(src image.Image, cell image.Rectangle) (r, g, b, a uint32) {
	width := cell.Max.X - cell.Min.X
	height := cell.Max.Y - cell.Min.Y
	if width*height == 0 {
		return src.At(cell.Min.X, cell.Min.Y).RGBA()
	}
	var red, green, blue, alpha uint32
	for x := cell.Min.X; x < cell.Max.X; x++ {
		for y := cell.Min.Y; y < cell.Max.Y; y++ {
			r, g, b, a := src.At(x, y).RGBA()
			red += r
			green += g
			blue += b
			alpha += a
		}
	}
	return red / uint32(width*height), green / uint32(width*height), blue / uint32(width*height), alpha / uint32(width*height)
}

// HorizontalAverageColorPicker picks horizontal-average color of center of given RectAngle area of src image.
type HorizontalAverageColorPicker struct{}

// Pick of HorizontalAverageColorPicker.
func (picker HorizontalAverageColorPicker) Pick(src image.Image, cell image.Rectangle) (r, g, b, a uint32) {
	var red, green, blue, alpha uint32
	width := cell.Max.X - cell.Min.X
	for x := cell.Min.X; x < cell.Max.X; x++ {
		r, g, b, a := src.At(x, (cell.Min.Y+cell.Max.Y)/2).RGBA()
		red += r
		green += g
		blue += b
		alpha += a
	}
	return red / uint32(width), green / uint32(width), blue / uint32(width), alpha / uint32(width)
}

// CenterColorPicker picks the very central point's color of given RectAngle area of src image.
type CenterColorPicker struct{}

// Pick of CenterColorPicker.
func (picker CenterColorPicker) Pick(src image.Image, cell image.Rectangle) (r, g, b, a uint32) {
	return src.At(int(float64(cell.Min.X+cell.Max.X)/2), int(float64(cell.Min.Y+cell.Max.Y)/2)).RGBA()
}

// LeftTopColorPicker picks color of left top (inital point) of given RectAngle of src image.
type LeftTopColorPicker struct{}

// Pick of LeftTopColorPicker.
func (picker LeftTopColorPicker) Pick(src image.Image, cell image.Rectangle) (r, g, b, a uint32) {
	return src.At(cell.Min.X, cell.Min.Y).RGBA()
}
