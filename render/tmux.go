package render

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

// tmuxQueryTimeout bounds how long we wait for tmux to answer an option query.
const tmuxQueryTimeout = time.Second

// inTmux reports whether the process runs inside a tmux pane.
func inTmux() bool {
	return os.Getenv("TMUX") != ""
}

// queryTmuxPassthrough returns the effective allow-passthrough value of the current pane.
// It is a variable so that tests can stub it.
var queryTmuxPassthrough = func() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tmuxQueryTimeout)
	defer cancel()
	args := []string{"display-message", "-p"}
	if pane := os.Getenv("TMUX_PANE"); pane != "" {
		args = append(args, "-t", pane)
	}
	out, err := exec.CommandContext(ctx, "tmux", append(args, "#{allow-passthrough}")...).Output()
	return strings.TrimSpace(string(out)), err
}

// tmuxPassthroughAllowed reports whether tmux forwards DCS passthrough sequences
// from the current pane to the outer terminal.
// tmux >= 3.3 has the allow-passthrough option (default off); older tmux has no
// such option, so the query succeeds with an empty value and passthrough is always allowed.
func tmuxPassthroughAllowed() bool {
	v, err := queryTmuxPassthrough()
	if err != nil {
		return false
	}
	switch v {
	case "on", "all", "":
		return true
	default:
		return false
	}
}

// tmuxPassthrough wraps an escape sequence so that tmux passes it through to the outer terminal.
// See https://github.com/tmux/tmux/wiki/FAQ#what-is-the-passthrough-escape-sequence-and-how-do-i-use-it
func tmuxPassthrough(seq string) string {
	return "\x1bPtmux;" + strings.ReplaceAll(seq, "\x1b", "\x1b\x1b") + "\x1b\\"
}
