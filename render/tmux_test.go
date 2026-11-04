package render

import (
	"bytes"
	"errors"
	"image/color"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Tests expect the behavior outside tmux unless they set TMUX explicitly.
	os.Unsetenv("TMUX")
	os.Exit(m.Run())
}

// tmuxEnv is a sample value of TMUX set by tmux in a pane.
const tmuxEnv = "/tmp/tmux-501/default,1,0"

// stubTmuxPassthrough replaces the tmux query for the duration of the test.
func stubTmuxPassthrough(t *testing.T, value string, err error) {
	t.Helper()
	orig := queryTmuxPassthrough
	queryTmuxPassthrough = func() (string, error) { return value, err }
	t.Cleanup(func() { queryTmuxPassthrough = orig })
}

func TestITermImageSupported_Tmux(t *testing.T) {
	queryErr := errors.New("tmux not found")
	tests := []struct {
		name        string
		tmux        string
		termProgram string
		lcTerminal  string
		passthrough string
		queryErr    error
		want        bool
	}{
		// AC-1
		{"tmux/LC_TERMINAL/on", tmuxEnv, "tmux", "iTerm2", "on", nil, true},
		{"tmux/LC_TERMINAL/all", tmuxEnv, "tmux", "iTerm2", "all", nil, true},
		{"tmux/LC_TERMINAL/old tmux", tmuxEnv, "tmux", "iTerm2", "", nil, true},
		{"old tmux/TERM_PROGRAM inherited/on", tmuxEnv, "iTerm.app", "", "on", nil, true},
		// AC-2
		{"tmux/LC_TERMINAL/off", tmuxEnv, "tmux", "iTerm2", "off", nil, false},
		{"tmux/TERM_PROGRAM inherited/off", tmuxEnv, "iTerm.app", "iTerm2", "off", nil, false},
		{"tmux/query failed", tmuxEnv, "tmux", "iTerm2", "", queryErr, false},
		{"tmux/unknown value", tmuxEnv, "tmux", "iTerm2", "maybe", nil, false},
		{"tmux/not iTerm2", tmuxEnv, "tmux", "", "on", nil, false},
		// AC-3
		{"no tmux/iTerm.app", "", "iTerm.app", "", "off", queryErr, true},
		{"no tmux/LC_TERMINAL only", "", "Apple_Terminal", "iTerm2", "on", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TMUX", tt.tmux)
			t.Setenv("TERM_PROGRAM", tt.termProgram)
			t.Setenv("LC_TERMINAL", tt.lcTerminal)
			stubTmuxPassthrough(t, tt.passthrough, tt.queryErr)
			if got := ITermImageSupported(); got != tt.want {
				t.Errorf("ITermImageSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}

// AC-4
func TestITerm_Render_Tmux(t *testing.T) {
	img := createTestImage(4, 4, color.RGBA{255, 0, 0, 255})

	t.Setenv("TMUX", "")
	plain := bytes.NewBuffer(nil)
	if err := (&ITerm{}).Render(plain, img); err != nil {
		t.Fatal(err)
	}
	seq := strings.TrimSuffix(plain.String(), "\n")
	if !strings.HasPrefix(seq, "\x1b]1337;File=;width=4px;height=4px;inline=1:") || !strings.HasSuffix(seq, "\a") {
		t.Fatalf("unexpected output outside tmux: %q", plain.String())
	}

	t.Setenv("TMUX", tmuxEnv)
	wrapped := bytes.NewBuffer(nil)
	if err := (&ITerm{}).Render(wrapped, img); err != nil {
		t.Fatal(err)
	}
	want := "\x1bPtmux;\x1b" + seq + "\x1b\\\n"
	if wrapped.String() != want {
		t.Errorf("ITerm.Render() in tmux\n got %q\nwant %q", wrapped.String(), want)
	}
}

func TestTmuxPassthrough(t *testing.T) {
	got := tmuxPassthrough("\x1b]1337;File=:AA\x1b\\")
	want := "\x1bPtmux;\x1b\x1b]1337;File=:AA\x1b\x1b\\\x1b\\"
	if got != want {
		t.Errorf("tmuxPassthrough() = %q, want %q", got, want)
	}
}
