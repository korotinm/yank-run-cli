package format

import (
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// IsStdoutTTY reports whether stdout is connected to a terminal.
func IsStdoutTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

// IsStderrTTY reports whether stderr is connected to a terminal.
func IsStderrTTY() bool {
	return isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())
}

// ConfigureColor disables color when not on a TTY or when forced off.
// NO_COLOR env var is honored automatically by fatih/color via its NoColor flag.
func ConfigureColor(forceOff bool) {
	if forceOff || !IsStdoutTTY() || os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}
}
