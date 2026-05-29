package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/korotinm/yank-run-cli/internal/format"
	"github.com/korotinm/yank-run-cli/internal/ids"
	"github.com/korotinm/yank-run-cli/internal/state"
)

var copyCmd = &cobra.Command{
	Use:   "copy REF",
	Short: "Copy snippet body to the system clipboard",
	Long: `Fetch a snippet and put its body on the system clipboard.

Requires a working clipboard backend on the host (pbcopy on macOS,
xclip/xsel/wl-copy on Linux, native on Windows).`,
	Args: cobra.ExactArgs(1),
	RunE: runCopy,
}

func init() {
	rootCmd.AddCommand(copyCmd)
}

func runCopy(cmd *cobra.Command, args []string) error {
	id, err := state.ResolveRef(args[0])
	if err != nil {
		return err
	}
	snip, err := newClient().Get(id)
	if err != nil {
		return err
	}
	if err := clipboard.WriteAll(snip.Body); err != nil {
		return fmt.Errorf("clipboard: %w", err)
	}
	if format.IsStderrTTY() {
		fmt.Fprintf(os.Stderr, "copied %s\n", ids.ShortPrefix(snip.ID))
	}
	return nil
}
