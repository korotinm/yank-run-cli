package cmd

import (
	"fmt"
	"os"
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

const defaultURL = "http://localhost:8080"

var rootCmd = &cobra.Command{
	Use:   "yank",
	Short: "yank.run command-line client",
	Long: `yank.run command-line client.

Search, store, and run code snippets without leaving your shell.

Configuration via env vars:
  YANK_URL      backend API base URL (default: http://localhost:8080)
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

// newClient builds the API client with current global flag values.
func newClient() *api.Client {
	return api.New(flagURL, flagTimeout, Version)
}
