package input

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

// ErrUsage signals a client-side validation/conflict error (exit code 2).
var ErrUsage = errors.New("usage error")

// ResolveBody picks exactly one body source from -b, -f, or piped stdin.
// Returns ErrUsage wrapped with a human message when sources conflict or none provided.
func ResolveBody(inline, file string) ([]byte, error) {
	stdinPiped := !isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd())

	switch {
	case inline != "" && file != "":
		return nil, fmt.Errorf("%w: use exactly one of -b, -f", ErrUsage)
	case (inline != "" || file != "") && stdinPiped:
		return nil, fmt.Errorf("%w: -b/-f conflicts with piped stdin", ErrUsage)
	case inline != "":
		return []byte(inline), nil
	case file != "":
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("read file %s: %w", file, err)
		}
		return data, nil
	case stdinPiped:
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("read stdin: %w", err)
		}
		return data, nil
	default:
		return nil, fmt.Errorf("%w: provide -b, -f, or pipe content into stdin", ErrUsage)
	}
}
