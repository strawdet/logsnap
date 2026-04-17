package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ryanfowler/logsnap/internal/snapshot"
	"github.com/spf13/cobra"
)

func init() {
	var dir string

	highlightCmd := &cobra.Command{
		Use:   "highlight",
		Short: "Manage highlights for a snapshot",
	}

	addCmd := &cobra.Command{
		Use:   "add <snapshot-id> <message>",
		Short: "Add a highlight message to a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.AddHighlight(dir, args[0], args[1]); err != nil {
				return fmt.Errorf("failed to add highlight: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Highlight added to snapshot %s\n", args[0])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <snapshot-id>",
		Short: "List highlights for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := snapshot.GetHighlights(dir, args[0])
			if err != nil {
				return fmt.Errorf("failed to get highlights: %w", err)
			}
			if len(h.Messages) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No highlights found.")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), strings.Join(h.Messages, "\n"))
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <snapshot-id> <message>",
		Short: "Remove a highlight message from a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.RemoveHighlight(dir, args[0], args[1]); err != nil {
				return fmt.Errorf("failed to remove highlight: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Highlight removed from snapshot %s\n", args[0])
			return nil
		},
	}

	defaultDir, _ := os.UserCacheDir()
	defaultDir = defaultDir + "/logsnap"

	for _, sub := range []*cobra.Command{addCmd, getCmd, removeCmd} {
		sub.Flags().StringVar(&dir, "dir", defaultDir, "Snapshot storage directory")
		highlightCmd.AddCommand(sub)
	}

	RootCmd.AddCommand(highlightCmd)
}
