package render

import (
	"bytes"
	"image"
	"image/color"
	"os"
	"testing"

	colorpkg "github.com/otiai10/gat/color"
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

func TestCellGrid_Render(t *testing.T) {
	img := createTestImage(100, 100, color.RGBA{255, 0, 0, 255})
	buf := bytes.NewBuffer(nil)

	grid := &CellGrid{
		Row:         10,
		Col:         20,
		Colorpicker: colorpkg.AverageColorPicker{},
		Placeholder: "  ",
	}

	err := grid.Render(buf, img)
	if err != nil {
		t.Fatalf("CellGrid.Render() error = %v", err)
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("CellGrid.Render() produced no output")
	}

	// Check that output contains ANSI escape codes
	if !bytes.Contains(buf.Bytes(), []byte("\x1b[48;5;")) {
		t.Error("CellGrid.Render() output does not contain expected ANSI escape codes")
	}
}

func TestCellGrid_Render_WithDefaults(t *testing.T) {
	img := createTestImage(100, 100, color.RGBA{0, 255, 0, 255})
	buf := bytes.NewBuffer(nil)

	// Test with nil Border and Colorpicker (should use defaults)
	grid := &CellGrid{
		Row: 10,
		Col: 20,
	}

	err := grid.Render(buf, img)
	if err != nil {
		t.Fatalf("CellGrid.Render() with defaults error = %v", err)
	}

	if buf.Len() == 0 {
		t.Error("CellGrid.Render() with defaults produced no output")
	}
}

func TestCellGrid_Render_TooSmall(t *testing.T) {
	img := createTestImage(100, 100, color.RGBA{0, 0, 255, 255})
	buf := bytes.NewBuffer(nil)

	grid := &CellGrid{
		Row: 1,
		Col: 1,
	}

	err := grid.Render(buf, img)
	if err == nil {
		t.Error("CellGrid.Render() expected error for too small canvas, got nil")
	}
}

func TestCellGrid_Fprint(t *testing.T) {
	t.Run("normal mode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		grid := &CellGrid{Placeholder: "  "}
		grid.Fprint(buf, 196) // red color code

		output := buf.String()
		if !bytes.Contains(buf.Bytes(), []byte("\x1b[48;5;196m")) {
			t.Errorf("Fprint() expected ANSI code for color 196, got %q", output)
		}
	})

	t.Run("debug mode", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		grid := &CellGrid{Placeholder: "  ", Debug: true}
		grid.Fprint(buf, 196)

		output := buf.String()
		if !bytes.Contains(buf.Bytes(), []byte("196")) {
			t.Errorf("Fprint() in debug mode expected color code in output, got %q", output)
		}
	})
}

func TestCellGrid_SetScale(t *testing.T) {
	grid := &CellGrid{}
	err := grid.SetScale(2.0)
	if err == nil {
		t.Error("CellGrid.SetScale() expected error (not supported), got nil")
	}
}

func TestITermImageSupported(t *testing.T) {
	// Save original env
	orig := os.Getenv("TERM_PROGRAM")
	defer os.Setenv("TERM_PROGRAM", orig)

	t.Run("iTerm.app", func(t *testing.T) {
		os.Setenv("TERM_PROGRAM", "iTerm.app")
		if !ITermImageSupported() {
			t.Error("ITermImageSupported() = false, want true for iTerm.app")
		}
	})

	t.Run("other terminal", func(t *testing.T) {
		os.Setenv("TERM_PROGRAM", "Terminal.app")
		if ITermImageSupported() {
			t.Error("ITermImageSupported() = true, want false for Terminal.app")
		}
	})

	t.Run("empty", func(t *testing.T) {
		os.Setenv("TERM_PROGRAM", "")
		if ITermImageSupported() {
			t.Error("ITermImageSupported() = true, want false for empty")
		}
	})
}

func TestGetDefaultRenderer(t *testing.T) {
	// Save original env
	orig := os.Getenv("TERM_PROGRAM")
	defer os.Setenv("TERM_PROGRAM", orig)

	t.Run("returns CellGrid for non-iTerm", func(t *testing.T) {
		os.Setenv("TERM_PROGRAM", "")
		r := GetDefaultRenderer()
		if _, ok := r.(*CellGrid); !ok {
			t.Errorf("GetDefaultRenderer() = %T, want *CellGrid", r)
		}
	})

	t.Run("returns ITerm for iTerm.app", func(t *testing.T) {
		os.Setenv("TERM_PROGRAM", "iTerm.app")
		r := GetDefaultRenderer()
		if _, ok := r.(*ITerm); !ok {
			t.Errorf("GetDefaultRenderer() = %T, want *ITerm", r)
		}
	})
}

func TestRenderer_Interface(t *testing.T) {
	// Verify all renderers implement the Renderer interface
	var _ Renderer = &CellGrid{}
	var _ Renderer = &ITerm{}
	var _ Renderer = &Sixel{}
}
