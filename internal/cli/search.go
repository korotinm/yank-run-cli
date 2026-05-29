package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/korotinm/yank-run-cli/internal/api"
	"github.com/korotinm/yank-run-cli/internal/format"
	"github.com/korotinm/yank-run-cli/internal/ids"
	"github.com/korotinm/yank-run-cli/internal/state"
)

var (
	searchLimit     int
	searchNoNumbers bool
)

var searchCmd = &cobra.Command{
	Use:   "search TERMS...",
	Short: "Full-text search snippets",
	Long: `Search snippets by full-text query. Terms are joined by spaces.

Results are cached to ~/.cache/yank/last-search.json so that follow-up
commands can refer to hits by their 0-based index (e.g. 'yank cat 0').`,
	Example: `  yank search vault kubectl
  yank search -n 5 docker
  yank search vault | fzf | awk '{print $1}' | xargs yank cat`,
	Args: cobra.MinimumNArgs(1),
	RunE: runSearch,
}

func init() {
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "n", 20, "max number of results (clamped server-side to 1..50)")
	searchCmd.Flags().BoolVar(&searchNoNumbers, "no-numbers", false, "suppress the # index column")
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	q := strings.Join(args, " ")
	hits, err := newClient().Search(q, searchLimit)
	if err != nil {
		return err
	}

	// Persist the result list for numeric-index resolution. Failure here
	// must not break the user's primary action — log and continue.
	if err := state.Save(q, hits); err != nil {
		fmt.Fprintf(os.Stderr, "yank: warn: failed to cache last search: %v\n", err)
	}

	if flagJSON {
		return json.NewEncoder(os.Stdout).Encode(hits)
	}

	if len(hits) == 0 {
		if format.IsStderrTTY() {
			fmt.Fprintln(os.Stderr, "(no results)")
		}
		return nil
	}

	printHits(hits)
	return nil
}

func printHits(hits []api.Hit) {
	tty := format.IsStdoutTTY()
	sep := "\t"
	if tty {
		sep = "  "
	}

	if tty && !searchNoNumbers {
		fmt.Printf("%-3s %-9s %s\n", "#", "ID", "PREVIEW")
	}

	for i, h := range hits {
		preview := h.BodyPreview
		var line string
		switch {
		case searchNoNumbers && tty:
			line = fmt.Sprintf("%-9s %s", ids.ShortPrefix(h.ID), preview)
		case searchNoNumbers:
			line = strings.Join([]string{ids.ShortPrefix(h.ID), preview}, sep)
		case tty:
			line = fmt.Sprintf("%-3d %-9s %s", i, ids.ShortPrefix(h.ID), preview)
		default:
			line = strings.Join([]string{fmt.Sprintf("%d", i), ids.ShortPrefix(h.ID), preview}, sep)
		}
		fmt.Println(line)
	}
}
