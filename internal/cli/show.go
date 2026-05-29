package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/korotinm/yank-run-cli/internal/api"
	"github.com/korotinm/yank-run-cli/internal/state"
)

var showCmd = &cobra.Command{
	Use:   "show REF",
	Short: "Human-readable snippet view",
	Long: `Display a snippet in a colored key-value layout followed by the body.

REF is either a 64-char lowercase hex id or a numeric index into the
last 'yank search' result list.`,
	Args: cobra.ExactArgs(1),
	RunE: runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	id, err := state.ResolveRef(args[0])
	if err != nil {
		return err
	}
	snip, err := newClient().Get(id)
	if err != nil {
		return err
	}

	if flagJSON {
		return json.NewEncoder(os.Stdout).Encode(snip)
	}

	printSnippet(snip)
	return nil
}

func printSnippet(s *api.Snippet) {
	label := color.New(color.FgCyan, color.Bold).SprintFunc()
	row := func(k, v string) {
		if v == "" {
			return
		}
		fmt.Printf("%-12s %s\n", label(k), v)
	}

	row("id", s.ID)
	row("title", s.Title)
	row("description", s.Description)
	if len(s.Tags) > 0 {
		row("tags", strings.Join(s.Tags, ", "))
	}
	row("author", s.Author)
	if s.CreatedAt != 0 {
		row("created", time.UnixMilli(s.CreatedAt).UTC().Format("2006-01-02 15:04:05 UTC"))
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 40))
	fmt.Print(s.Body)
	if !strings.HasSuffix(s.Body, "\n") {
		fmt.Println()
	}
}
