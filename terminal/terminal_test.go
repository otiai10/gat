package terminal

import "testing"

func TestGetTerminal(t *testing.T) {
	// Under `go test`, stdout/stdin are usually not terminals, so the fallback is used.
	// In either case, the size must be non-zero and must not panic.
	rect := GetTerminal()
	if rect.Row == 0 || rect.Col == 0 {
		t.Errorf("GetTerminal() = %+v, want non-zero size", rect)
	}
}
