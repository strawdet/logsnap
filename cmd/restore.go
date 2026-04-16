package cmd

import (
	"fmt"
	"os"

	"github.com/nicholasgasior/logsnap/internal/snapshot"
	"github.com/spf13/cobra"
)

var (
	restoreOutput string
)

func init() {
	restoreCmd := &cobra.Command{
		Use:   "restore <snapshot-id>",
		Short: "Restore log entries from a snapshot to a file or stdout",
		Args:  cobra.ExactArgs(1),
		RunE:  runRestore,
	}

	restoreCmd.Flags().StringVarP(&restoreOutput, "output", "o", "-", "Output file path (use '-' for stdout)")
	restoreCmd.Flags().StringVar(&snapshotDir, "dir", defaultSnapshotDir(), "Directory storing snapshots")

	rootCmd.AddCommand(restoreCmd)
}

func runRestore(cmd *cobra.Command, args []string) error {
	snapshotID := args[0]

	result, err := snapshot.Restore(snapshotDir, snapshotID, restoreOutput)
	if err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}

	if restoreOutput != "-" {
		fmt.Fprintf(os.Stderr, "Restored %d entries from snapshot %s to %s\n",
			result.EntryCount, result.SnapshotID, result.OutputPath)
	}
	return nil
}
