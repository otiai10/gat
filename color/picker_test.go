package color

import (
	"image"
	"image/color"
	"testing"
)

// createTestImage creates a simple test image with a solid color
func createTestImage(width, height int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

// createGradientImage creates a gradient test image
func createGradientImage(width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gray := uint8((x + y) * 255 / (width + height))
			img.Set(x, y, color.RGBA{gray, gray, gray, 255})
		}
	}
	return img
}

func TestAverageColorPicker_Pick(t *testing.T) {
	picker := AverageColorPicker{}

	t.Run("solid red image", func(t *testing.T) {
		img := createTestImage(10, 10, color.RGBA{255, 0, 0, 255})
		r, g, b, a := picker.Pick(img, image.Rect(0, 0, 10, 10))
		if r>>8 != 255 || g != 0 || b != 0 {
			t.Errorf("Expected red, got r=%d, g=%d, b=%d, a=%d", r>>8, g>>8, b>>8, a>>8)
		}
	})

	t.Run("zero area falls back to point", func(t *testing.T) {
		img := createTestImage(10, 10, color.RGBA{0, 255, 0, 255})
		r, g, b, _ := picker.Pick(img, image.Rect(5, 5, 5, 5))
		if g>>8 != 255 {
			t.Errorf("Expected green from point, got r=%d, g=%d, b=%d", r>>8, g>>8, b>>8)
		}
	})
}

func TestHorizontalAverageColorPicker_Pick(t *testing.T) {
	picker := HorizontalAverageColorPicker{}

	t.Run("solid blue image", func(t *testing.T) {
		img := createTestImage(10, 10, color.RGBA{0, 0, 255, 255})
		r, g, b, _ := picker.Pick(img, image.Rect(0, 0, 10, 10))
		if b>>8 != 255 || r != 0 || g != 0 {
			t.Errorf("Expected blue, got r=%d, g=%d, b=%d", r>>8, g>>8, b>>8)
		}
	})
}

func TestCenterColorPicker_Pick(t *testing.T) {
	picker := CenterColorPicker{}

	t.Run("picks center pixel", func(t *testing.T) {
		img := createTestImage(10, 10, color.RGBA{128, 128, 128, 255})
		r, g, b, _ := picker.Pick(img, image.Rect(0, 0, 10, 10))
		if r>>8 != 128 || g>>8 != 128 || b>>8 != 128 {
			t.Errorf("Expected gray (128), got r=%d, g=%d, b=%d", r>>8, g>>8, b>>8)
		}
	})
}

func TestLeftTopColorPicker_Pick(t *testing.T) {
	picker := LeftTopColorPicker{}

	t.Run("picks top-left pixel", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		// Set only top-left pixel to red
		img.Set(0, 0, color.RGBA{255, 0, 0, 255})
		r, g, b, _ := picker.Pick(img, image.Rect(0, 0, 10, 10))
		if r>>8 != 255 || g != 0 || b != 0 {
			t.Errorf("Expected red from top-left, got r=%d, g=%d, b=%d", r>>8, g>>8, b>>8)
		}
	})
}

func TestPicker_Interface(t *testing.T) {
	// Verify all pickers implement the Picker interface
	var _ Picker = AverageColorPicker{}
	var _ Picker = HorizontalAverageColorPicker{}
	var _ Picker = CenterColorPicker{}
	var _ Picker = LeftTopColorPicker{}
}
