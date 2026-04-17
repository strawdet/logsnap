package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	var dir string

	noteCmd := &cobra.Command{
		Use:   "note",
		Short: "Manage notes attached to snapshots",
	}

	addCmd := &cobra.Command{
		Use:   "add <snapshot-id> <text>",
		Short: "Add or update a note on a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.AddNote(dir, args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Note saved for snapshot %s\n", args[0])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <snapshot-id>",
		Short: "Get the note for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			note, err := snapshot.GetNote(dir, args[0])
			if err != nil {
				return err
			}
			fmt.Printf("Note: %s\nCreated: %s\nUpdated: %s\n",
				note.Text,
				note.CreatedAt.Format("2006-01-02 15:04:05"),
				note.UpdatedAt.Format("2006-01-02 15:04:05"),
			)
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <snapshot-id>",
		Short: "Remove the note from a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.RemoveNote(dir, args[0]); err != nil {
				return err
			}
			fmt.Printf("Note removed from snapshot %s\n", args[0])
			return nil
		},
	}

	for _, sub := range []*cobra.Command{addCmd, getCmd, removeCmd} {
		sub.Flags().StringVar(&dir, "dir", "snapshots", "Snapshot storage directory")
		noteCmd.AddCommand(sub)
	}

	if rootCmd := getRootCmd(); rootCmd != nil {
		rootCmd.AddCommand(noteCmd)
	} else {
		fmt.Fprintln(os.Stderr, "warn: rootCmd not available")
	}
}
