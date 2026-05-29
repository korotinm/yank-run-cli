package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/korotinm/yank-run-cli/internal/state"
)

var catCmd = &cobra.Command{
	Use:   "cat REF",
	Short: "Print a snippet body to stdout",
	Long: `Print the raw snippet body to stdout. Designed for shell composition:

  yank cat <id> | bash
  yank cat 0 > script.sh

REF is either a 64-char lowercase hex id or a numeric index (0..N) into
the last 'yank search' result list. No trailing newline is added by the
CLI — what the server stored is what you get.`,
	Args: cobra.ExactArgs(1),
	RunE: runCat,
}

func init() {
	rootCmd.AddCommand(catCmd)
}

func runCat(cmd *cobra.Command, args []string) error {
	id, err := state.ResolveRef(args[0])
	if err != nil {
		return err
	}
	snip, err := newClient().Get(id)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprint(os.Stdout, snip.Body); err != nil {
		return err
	}
	return nil
}
