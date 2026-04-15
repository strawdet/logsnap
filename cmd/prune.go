package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

var (
	pruneKeepLast int
	pruneOlderDays int
	pruneDryRun bool
)

func init() {
	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove old snapshots based on age or count",
		RunE:  runPrune,
	}

	pruneCmd.Flags().IntVar(&pruneKeepLast, "keep-last", 0, "Keep the N most recent snapshots (0 = no limit)")
	pruneCmd.Flags().IntVar(&pruneOlderDays, "older-than", 0, "Remove snapshots older than N days (0 = no limit)")
	pruneCmd.Flags().BoolVar(&pruneDryRun, "dry-run", false, "List snapshots that would be removed without deleting them")

	rootCmd.AddCommand(pruneCmd)
}

func runPrune(cmd *cobra.Command, args []string) error {
	dir := os.Getenv("LOGSNAP_DIR")
	if dir == "" {
		dir = ".logsnap"
	}

	opts := snapshot.PruneOptions{
		KeepLast: pruneKeepLast,
		DryRun:   pruneDryRun,
	}
	if pruneOlderDays > 0 {
		opts.OlderThan = time.Now().AddDate(0, 0, -pruneOlderDays)
	}

	if opts.KeepLast == 0 && opts.OlderThan.IsZero() {
		return fmt.Errorf("specify at least one of --keep-last or --older-than")
	}

	result, err := snapshot.Prune(dir, opts)
	if err != nil {
		return fmt.Errorf("prune failed: %w", err)
	}

	if pruneDryRun {
		fmt.Fprintf(cmd.OutOrStdout(), "Dry run — would remove %d snapshot(s):\n", len(result.Removed))
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Removed %d snapshot(s):\n", len(result.Removed))
	}
	for _, id := range result.Removed {
		fmt.Fprintf(cmd.OutOrStdout(), "  - %s\n", id)
	}
	if len(result.Removed) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "  (none)")
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Kept %d snapshot(s).\n", len(result.Kept))
	return nil
}
