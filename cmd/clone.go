package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	var label string

	cloneCmd := &cobra.Command{
		Use:   "clone <snapshot-id>",
		Short: "Clone an existing snapshot into a new one",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runClone(args[0], label)
		},
	}

	cloneCmd.Flags().StringVar(&label, "label", "", "Label for the cloned snapshot (default: original label + ' (clone)'")

	rootCmd.AddCommand(cloneCmd)
}

func runClone(id, label string) error {
	dir := os.Getenv("LOGSNAP_DIR")
	if dir == "" {
		dir = ".logsnap"
	}

	cloned, err := snapshot.CloneSnapshot(dir, id, label)
	if err != nil {
		return fmt.Errorf("clone failed: %w", err)
	}

	fmt.Printf("Cloned snapshot %q → new ID: %s\n", id, cloned.ID)
	fmt.Printf("Label: %s\n", cloned.Label)
	fmt.Printf("Entries: %d\n", len(cloned.Entries))
	return nil
}
