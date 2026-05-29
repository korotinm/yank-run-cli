package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/korotinm/yank-run-cli/internal/format"
	"github.com/korotinm/yank-run-cli/internal/state"
)

var openCmd = &cobra.Command{
	Use:   "open REF",
	Short: "Open the snippet URL in the default browser",
	Long: `Open a snippet URL in the user's default browser.

URL resolution:
  1. If YANK_WEB_URL is set: open ${YANK_WEB_URL}/${id}
  2. Otherwise: open ${YANK_URL}/api/snippets/${id} (raw JSON)`,
	Args: cobra.ExactArgs(1),
	RunE: runOpen,
}

func init() {
	rootCmd.AddCommand(openCmd)
}

func runOpen(cmd *cobra.Command, args []string) error {
	id, err := state.ResolveRef(args[0])
	if err != nil {
		return err
	}

	var url string
	if web := strings.TrimRight(os.Getenv("YANK_WEB_URL"), "/"); web != "" {
		url = fmt.Sprintf("%s/%s", web, id)
	} else {
		url = fmt.Sprintf("%s/api/snippets/%s", strings.TrimRight(flagURL, "/"), id)
	}

	if err := browser.OpenURL(url); err != nil {
		return fmt.Errorf("open browser: %w", err)
	}
	if format.IsStderrTTY() {
		fmt.Fprintf(os.Stderr, "opening %s\n", url)
	}
	return nil
}
