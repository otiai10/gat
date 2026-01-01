package color

import (
	"bytes"
	"testing"
)

func TestRGB2Hex(t *testing.T) {
	tests := []struct {
		name     string
		r, g, b  uint32
		expected uint32
	}{
		{"black", 0, 0, 0, 0x000000},
		{"white", 0xffff, 0xffff, 0xffff, 0xffffff},
		{"red", 0xffff, 0, 0, 0xff0000},
		{"green", 0, 0xffff, 0, 0x00ff00},
		{"blue", 0, 0, 0xffff, 0x0000ff},
		{"mid gray", 0x8080, 0x8080, 0x8080, 0x808080},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RGB2Hex(tt.r, tt.g, tt.b)
			if got != tt.expected {
				t.Errorf("RGB2Hex(%d, %d, %d) = 0x%06x, want 0x%06x", tt.r, tt.g, tt.b, got, tt.expected)
			}
		})
	}
}

func TestGetCodeByRGBA(t *testing.T) {
	tests := []struct {
		name        string
		r, g, b, a  uint32
		expectedMin int
		expectedMax int
	}{
		{"black", 0, 0, 0, 0xffff, 0, 0},
		{"white", 0xffff, 0xffff, 0xffff, 0xffff, 15, 15},
		{"red", 0xffff, 0, 0, 0xffff, 9, 9},
		{"green", 0, 0xffff, 0, 0xffff, 10, 10},
		{"blue", 0, 0, 0xffff, 0xffff, 12, 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCodeByRGBA(tt.r, tt.g, tt.b, tt.a)
			if got < tt.expectedMin || got > tt.expectedMax {
				t.Errorf("GetCodeByRGBA(%d, %d, %d, %d) = %d, want between %d and %d",
					tt.r, tt.g, tt.b, tt.a, got, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

func TestFindApproximateColorCode(t *testing.T) {
	// Test that approximate colors return valid codes (0-255)
	tests := []struct {
		name    string
		r, g, b uint32
	}{
		{"near black", 0x0100, 0x0100, 0x0100},
		{"near white", 0xfe00, 0xfe00, 0xfe00},
		{"arbitrary", 0x1234, 0x5678, 0x9abc},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindApproximateColorCode(tt.r, tt.g, tt.b)
			if got < 0 || got > 255 {
				t.Errorf("FindApproximateColorCode(%d, %d, %d) = %d, want 0-255",
					tt.r, tt.g, tt.b, got)
			}
		})
	}
}

func TestDump(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	Dump(buf)

	// Dump should output 256 color codes
	output := buf.String()
	if len(output) == 0 {
		t.Error("Dump() produced no output")
	}
	// Check that output contains ANSI escape codes
	if !bytes.Contains(buf.Bytes(), []byte("\x1b[48;5;")) {
		t.Error("Dump() output does not contain expected ANSI escape codes")
	}
}
