package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

var (
	replayLevel string
	replayDelay int
)

func init() {
	replayCmd := &cobra.Command{
		Use:   "replay <snapshot-id>",
		Short: "Replay log entries from a snapshot to stdout",
		Args:  cobra.ExactArgs(1),
		RunE:  runReplay,
	}

	replayCmd.Flags().StringVar(&replayLevel, "level", "", "filter entries by log level (e.g. error, info)")
	replayCmd.Flags().IntVar(&replayDelay, "delay", 0, "delay in milliseconds between entries")

	rootCmd.AddCommand(replayCmd)
}

func runReplay(cmd *cobra.Command, args []string) error {
	snapshotID := args[0]
	dir := os.Getenv("LOGSNAP_DIR")
	if dir == "" {
		dir = ".logsnap"
	}

	// Resolve tag if needed
	resolved, err := snapshot.ResolveTag(dir, snapshotID)
	if err == nil {
		snapshotID = resolved
	}

	snap, err := snapshot.Load(dir, snapshotID)
	if err != nil {
		return fmt.Errorf("failed to load snapshot %q: %w", snapshotID, err)
	}

	opts := snapshot.ReplayOptions{
		Filter: replayLevel,
		Delay:  time.Duration(replayDelay) * time.Millisecond,
		Writer: cmd.OutOrStdout(),
	}

	result, err := snapshot.Replay(snap, opts)
	if err != nil {
		return fmt.Errorf("replay error: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\nReplayed %d entries, skipped %d.\n", result.Replayed, result.Skipped)
	return nil
}
