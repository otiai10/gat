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
	var _ Renderer = &Kitty{}
}

func TestKitty_Render(t *testing.T) {
	img := createTestImage(100, 100, color.RGBA{255, 0, 0, 255})
	buf := bytes.NewBuffer(nil)

	kitty := &Kitty{Scale: 1}
	err := kitty.Render(buf, img)
	if err != nil {
		t.Fatalf("Kitty.Render() error = %v", err)
	}

	output := buf.Bytes()
	if len(output) == 0 {
		t.Error("Kitty.Render() produced no output")
	}

	// Check that output contains Kitty Graphics Protocol escape sequence
	// Format: <ESC>_G...
	if !bytes.Contains(output, []byte("\x1b_G")) {
		t.Error("Kitty.Render() output does not contain expected escape sequence \\x1b_G")
	}

	// Check that output contains the action and format parameters
	if !bytes.Contains(output, []byte("a=T")) {
		t.Error("Kitty.Render() output does not contain action parameter a=T")
	}
	if !bytes.Contains(output, []byte("f=100")) {
		t.Error("Kitty.Render() output does not contain format parameter f=100")
	}

	// Check that output ends with the terminator
	if !bytes.Contains(output, []byte("\x1b\\")) {
		t.Error("Kitty.Render() output does not contain terminator \\x1b\\")
	}
}

func TestKitty_Render_LargeImage(t *testing.T) {
	// Create a larger image with varied colors to ensure PNG doesn't compress too much
	// This ensures the base64 output exceeds 4096 bytes and requires chunking
	img := image.NewRGBA(image.Rect(0, 0, 800, 800))
	for y := 0; y < 800; y++ {
		for x := 0; x < 800; x++ {
			// Create a gradient pattern that doesn't compress well
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), uint8((x + y) % 256), 255})
		}
	}
	buf := bytes.NewBuffer(nil)

	kitty := &Kitty{Scale: 1}
	err := kitty.Render(buf, img)
	if err != nil {
		t.Fatalf("Kitty.Render() error = %v", err)
	}

	output := buf.Bytes()
	// For large images, we should see m=1 (more chunks) followed by m=0 (last chunk)
	if !bytes.Contains(output, []byte("m=1")) {
		t.Error("Kitty.Render() large image should have m=1 for chunked transfer")
	}
	if !bytes.Contains(output, []byte("m=0")) {
		t.Error("Kitty.Render() large image should have m=0 for final chunk")
	}
}

func TestKitty_SetScale(t *testing.T) {
	kitty := &Kitty{}
	err := kitty.SetScale(2.0)
	if err != nil {
		t.Errorf("Kitty.SetScale() unexpected error = %v", err)
	}
	if kitty.Scale != 2.0 {
		t.Errorf("Kitty.SetScale() Scale = %v, want 2.0", kitty.Scale)
	}
}

func TestKittySupported(t *testing.T) {
	// Save original env
	origKittyWindowID := os.Getenv("KITTY_WINDOW_ID")
	origTermProgram := os.Getenv("TERM_PROGRAM")
	defer func() {
		os.Setenv("KITTY_WINDOW_ID", origKittyWindowID)
		os.Setenv("TERM_PROGRAM", origTermProgram)
	}()

	t.Run("Kitty terminal", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "1")
		os.Setenv("TERM_PROGRAM", "")
		if !KittySupported() {
			t.Error("KittySupported() = false, want true for Kitty terminal")
		}
	})

	t.Run("Ghostty terminal", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "")
		os.Setenv("TERM_PROGRAM", "ghostty")
		if !KittySupported() {
			t.Error("KittySupported() = false, want true for Ghostty")
		}
	})

	t.Run("WezTerm terminal", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "")
		os.Setenv("TERM_PROGRAM", "WezTerm")
		if !KittySupported() {
			t.Error("KittySupported() = false, want true for WezTerm")
		}
	})

	t.Run("unsupported terminal", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "")
		os.Setenv("TERM_PROGRAM", "Terminal.app")
		if KittySupported() {
			t.Error("KittySupported() = true, want false for Terminal.app")
		}
	})

	t.Run("empty env", func(t *testing.T) {
		os.Setenv("KITTY_WINDOW_ID", "")
		os.Setenv("TERM_PROGRAM", "")
		if KittySupported() {
			t.Error("KittySupported() = true, want false for empty env")
		}
	})
}
