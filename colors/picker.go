package colors

import "image"

// Picker ...
type Picker func(image.Image, image.Rectangle) (r, g, b, a uint32)

// AverageColor ...
func AverageColor(src image.Image, cell image.Rectangle) (r, g, b, a uint32) {
	var red, green, blue, alpha uint32
	width := cell.Max.X - cell.Min.X
	height := cell.Max.Y - cell.Min.Y
	for x := cell.Min.X; x < cell.Max.X; x++ {
		for y := cell.Min.Y; y < cell.Max.Y; y++ {
			r, g, b, a := src.At(x, y).RGBA()
			red += r
			green += g
			blue += b
			alpha += a
		}
	}
	// log.Println(red/uint32(width), green/uint32(width), blue/uint32(width))
	return red / uint32(width*height), green / uint32(width*height), blue / uint32(width*height), alpha / uint32(width*height)
}

// AverageColorX ...
func AverageColorX(src image.Image, cell image.Rectangle) (r, g, b, a uint32) {
	var red, green, blue, alpha uint32
	width := cell.Max.X - cell.Min.X
	for x := cell.Min.X; x < cell.Max.X; x++ {
		r, g, b, a := src.At(x, (cell.Min.Y+cell.Max.Y)/2).RGBA()
		red += r
		green += g
		blue += b
		alpha += a
	}
	// log.Println(red/uint32(width), green/uint32(width), blue/uint32(width))
	return red / uint32(width), green / uint32(width), blue / uint32(width), alpha / uint32(width)
}

// CenterColor ...
func CenterColor(src image.Image, cell image.Rectangle) (r, g, b, a uint32) {
	return src.At((cell.Min.X+cell.Max.X)/2, (cell.Min.Y+cell.Max.Y)/2).RGBA()
}
