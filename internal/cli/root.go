package cli

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"
	"github.com/korotinm/yank-run-cli/internal/api"
	"github.com/korotinm/yank-run-cli/internal/format"
)

// Build-time injected (see Makefile).
var (
	Version = "dev"
	Commit  = "none"
)

// Global flag values.
var (
	flagURL     string
	flagTimeout time.Duration
	flagJSON    bool
	flagNoColor bool
)

const defaultURL = "https://api.yank.run"

var rootCmd = &cobra.Command{
	Use:   "yank",
	Short: "yank.run command-line client",
	Long: `yank.run command-line client.

Search, store, and run code snippets without leaving your shell.

Configuration via env vars:
  YANK_URL      backend API base URL (default: https://api.yank.run)
  YANK_WEB_URL  frontend URL used by 'yank open'
  NO_COLOR      disable ANSI colors when set (any value)`,
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		format.ConfigureColor(flagNoColor)
	},
}

// Execute runs the root command, mapping returned errors to exit codes.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(format.HandleError(err))
	}
}

func init() {
	resolveVersion()

	urlDefault := os.Getenv("YANK_URL")
	if urlDefault == "" {
		urlDefault = defaultURL
	}

	rootCmd.PersistentFlags().StringVar(&flagURL, "url", urlDefault, "backend API base URL")
	rootCmd.PersistentFlags().DurationVar(&flagTimeout, "timeout", 10*time.Second, "HTTP request timeout")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "emit JSON instead of human/plain output")
	rootCmd.PersistentFlags().BoolVar(&flagNoColor, "no-color", false, "disable ANSI colors")

	rootCmd.Version = fmt.Sprintf("%s (commit %s)", Version, Commit)
	rootCmd.SetVersionTemplate("yank {{.Version}}\n")
}

// resolveVersion fills Version/Commit from runtime/debug.BuildInfo when the
// Makefile's ldflags weren't used (e.g. `go install <module>@<ver>`). When the
// linker has already injected real values, this function is a no-op.
func resolveVersion() {
	if Version != "dev" && Commit != "none" {
		return
	}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	// Module version: a tag like "v0.1.0", a pseudo-version, or "(devel)"
	// when built from a local checkout without go install.
	if Version == "dev" && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		Version = bi.Main.Version
	}
	// VCS settings are present for `go install module@ver` and for `go build`
	// from inside a git tree (Go ≥1.18). Absent for `go run`.
	var (
		rev      string
		modified bool
	)
	for _, s := range bi.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			modified = s.Value == "true"
		}
	}
	if Commit == "none" && rev != "" {
		short := rev
		if len(short) > 7 {
			short = short[:7]
		}
		if modified {
			short += "-dirty"
		}
		Commit = short
	}
	// If we still have no Version but got a revision, surface it instead of "dev".
	if Version == "dev" && rev != "" {
		Version = Commit
	}
}

// newClient builds the API client with current global flag values.
func newClient() *api.Client {
	return api.New(flagURL, flagTimeout, Version)
}
