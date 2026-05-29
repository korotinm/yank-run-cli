package format

import (
	"errors"
	"fmt"
	"os"

	"github.com/korotinm/yank-run-cli/internal/api"
	"github.com/korotinm/yank-run-cli/internal/input"
	"github.com/korotinm/yank-run-cli/internal/state"
)

// Exit codes per CLI_PLAN.md.
const (
	ExitOK            = 0
	ExitUsage         = 2 // client-side validation / conflicting flags
	ExitBadRequest    = 3 // server 4xx (non-429), bad REF, not found
	ExitRateLimited   = 4 // server 429
	ExitServerNetwork = 5 // 5xx, network failures, clipboard, etc.
)

// HandleError maps an error to (message, exit code) following the conventions.
// It writes the message to stderr and returns the code; the caller should os.Exit it.
func HandleError(err error) int {
	if err == nil {
		return ExitOK
	}

	switch {
	case errors.Is(err, input.ErrUsage):
		fmt.Fprintln(os.Stderr, "yank: "+err.Error())
		return ExitUsage
	case errors.Is(err, state.ErrNoRecentSearch),
		errors.Is(err, state.ErrBadRef):
		fmt.Fprintln(os.Stderr, "yank: "+err.Error())
		return ExitBadRequest
	}

	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		fmt.Fprintln(os.Stderr, "yank: "+apiErr.Message)
		switch {
		case apiErr.Status == 429:
			return ExitRateLimited
		case apiErr.Status >= 500:
			return ExitServerNetwork
		case apiErr.Status >= 400:
			return ExitBadRequest
		}
	}

	fmt.Fprintln(os.Stderr, "yank: "+err.Error())
	return ExitServerNetwork
}
