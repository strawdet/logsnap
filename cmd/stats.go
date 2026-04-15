package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	statsCmd := &cobra.Command{
		Use:   "stats <snapshot-id>",
		Short: "Display aggregate statistics for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE:  runStats,
	}
	rootCmd.AddCommand(statsCmd)
}

func runStats(cmd *cobra.Command, args []string) error {
	snapshotID := args[0]
	dir, err := snapshotDir()
	if err != nil {
		return err
	}

	// Resolve tag if needed
	resolved, err := snapshot.ResolveTag(dir, snapshotID)
	if err == nil {
		snapshotID = resolved
	}

	snap, err := snapshot.Load(dir, snapshotID)
	if err != nil {
		return fmt.Errorf("failed to load snapshot: %w", err)
	}

	stats, err := snapshot.ComputeStats(snap)
	if err != nil {
		return fmt.Errorf("failed to compute stats: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Snapshot ID:\t%s\n", stats.SnapshotID)
	fmt.Fprintf(w, "Label:\t%s\n", stats.Label)
	fmt.Fprintf(w, "Total Entries:\t%d\n", stats.TotalCount)
	w.Flush()

	fmt.Println("\nEntries by Level:")
	levels := make([]string, 0, len(stats.LevelCounts))
	for l := range stats.LevelCounts {
		levels = append(levels, l)
	}
	sort.Strings(levels)
	for _, l := range levels {
		fmt.Fprintf(w, "  %s:\t%d\n", l, stats.LevelCounts[l])
	}
	w.Flush()

	if len(stats.TopMessages) > 0 {
		fmt.Println("\nTop Messages:")
		for i, msg := range stats.TopMessages {
			fmt.Printf("  %d. %s\n", i+1, msg)
		}
	}

	return nil
}
