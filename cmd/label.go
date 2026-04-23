package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	var dir string

	labelCmd := &cobra.Command{
		Use:   "label",
		Short: "Manage labels on snapshots",
	}

	addCmd := &cobra.Command{
		Use:   "add <snapshot-id> <label>",
		Short: "Add a label to a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.AddLabel(dir, args[0], args[1]); err != nil {
				return fmt.Errorf("add label: %w", err)
			}
			fmt.Printf("Label %q added to snapshot %s\n", args[1], args[0])
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <snapshot-id> <label>",
		Short: "Remove a label from a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.RemoveLabel(dir, args[0], args[1]); err != nil {
				return fmt.Errorf("remove label: %w", err)
			}
			fmt.Printf("Label %q removed from snapshot %s\n", args[1], args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <snapshot-id>",
		Short: "List all labels for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			labels, err := snapshot.ListLabelsForSnapshot(dir, args[0])
			if err != nil {
				return fmt.Errorf("list labels: %w", err)
			}
			if len(labels) == 0 {
				fmt.Println("No labels found.")
				return nil
			}
			fmt.Println(strings.Join(labels, "\n"))
			return nil
		},
	}

	byLabelCmd := &cobra.Command{
		Use:   "snapshots <label>",
		Short: "List all snapshots with a given label",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, err := snapshot.GetSnapshotsByLabel(dir, args[0])
			if err != nil {
				return fmt.Errorf("get snapshots by label: %w", err)
			}
			if len(ids) == 0 {
				fmt.Printf("No snapshots with label %q.\n", args[0])
				return nil
			}
			for _, id := range ids {
				fmt.Println(id)
			}
			return nil
		},
	}

	for _, sub := range []*cobra.Command{addCmd, removeCmd, listCmd, byLabelCmd} {
		sub.Flags().StringVar(&dir, "dir", defaultSnapshotDir(), "snapshot storage directory")
		labelCmd.AddCommand(sub)
	}

	rootCmd.AddCommand(labelCmd)
}

func defaultSnapshotDir() string {
	if d := os.Getenv("LOGSNAP_DIR"); d != "" {
		return d
	}
	return ".logsnap"
}
