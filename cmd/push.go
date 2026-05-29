package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/korotinm/yank-run-cli/internal/api"
	"github.com/korotinm/yank-run-cli/internal/input"
)

var (
	pushBody        string
	pushFile        string
	pushTitle       string
	pushDescription string
	pushTags        []string
	pushAuthor      string
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Create a snippet (body from -b, -f, or stdin)",
	Long: `Create a snippet on the backend.

Body source: exactly one of -b, -f, or piped stdin.
The backend is content-addressed — a duplicate body returns the existing
id with HTTP 200; in that case the CLI prints "(existing)" to stderr.`,
	Example: `  yank push -b 'kubectl get pods -A' -t "List pods" --tag k8s
  yank push -f deploy/manifest.yaml --tag k8s,prod
  history | tail -1 | yank push --tag bash`,
	RunE: runPush,
}

func init() {
	pushCmd.Flags().StringVarP(&pushBody, "body", "b", "", "inline body string")
	pushCmd.Flags().StringVarP(&pushFile, "file", "f", "", "read body from file")
	pushCmd.Flags().StringVarP(&pushTitle, "title", "t", "", "snippet title (max 40 chars)")
	pushCmd.Flags().StringVarP(&pushDescription, "description", "d", "", "description (max 500 bytes, indexed)")
	pushCmd.Flags().StringSliceVar(&pushTags, "tag", nil, "tag (repeatable or comma-separated, max 5)")
	pushCmd.Flags().StringVar(&pushAuthor, "author", "", "author handle")
	rootCmd.AddCommand(pushCmd)
}

func runPush(cmd *cobra.Command, args []string) error {
	body, err := input.ResolveBody(pushBody, pushFile)
	if err != nil {
		return err
	}

	req := api.CreateReq{
		Body:        string(body),
		Title:       strings.TrimSpace(pushTitle),
		Description: strings.TrimSpace(pushDescription),
		Tags:        normalizeTags(pushTags),
		Author:      strings.TrimSpace(pushAuthor),
	}

	resp, err := newClient().Create(req)
	if err != nil {
		return err
	}

	if flagJSON {
		return json.NewEncoder(os.Stdout).Encode(resp)
	}

	fmt.Printf("%s\t%s\n", resp.ID, resp.URL)
	if !resp.New {
		fmt.Fprintln(os.Stderr, "(existing)")
	}
	return nil
}

// normalizeTags drops empties; the cobra StringSlice flag already splits commas.
func normalizeTags(in []string) []string {
	out := make([]string, 0, len(in))
	for _, t := range in {
		t = strings.TrimSpace(t)
		if t != "" {
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
