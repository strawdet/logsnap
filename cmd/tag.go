package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

var snapshotDir string

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage snapshot tags",
	}

	addTagCmd := &cobra.Command{
		Use:   "add <tag> <snapshot-id>",
		Short: "Associate a tag with a snapshot ID",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tag, id := args[0], args[1]
			if err := snapshot.TagSnapshot(snapshotDir, tag, id); err != nil {
				return fmt.Errorf("tagging snapshot: %w", err)
			}
			fmt.Printf("Tagged snapshot %s as %q\n", id, tag)
			return nil
		},
	}

	removeTagCmd := &cobra.Command{
		Use:   "remove <tag>",
		Short: "Remove a tag",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.RemoveTag(snapshotDir, args[0]); err != nil {
				return fmt.Errorf("removing tag: %w", err)
			}
			fmt.Printf("Removed tag %q\n", args[0])
			return nil
		},
	}

	listTagsCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			index, err := snapshot.LoadTagIndex(snapshotDir)
			if err != nil {
				return fmt.Errorf("loading tags: %w", err)
			}
			if len(index) == 0 {
				fmt.Println("No tags defined.")
				return nil
			}
			for tag, id := range index {
				fmt.Fprintf(os.Stdout, "%-20s %s\n", tag, id)
			}
			return nil
		},
	}

	tagCmd.PersistentFlags().StringVar(&snapshotDir, "dir", "snapshots", "Directory where snapshots are stored")
	tagCmd.AddCommand(addTagCmd, removeTagCmd, listTagsCmd)
	rootCmd.AddCommand(tagCmd)
}
