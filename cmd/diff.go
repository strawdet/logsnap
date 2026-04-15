package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/logsnap/internal/diff"
	"github.com/user/logsnap/internal/snapshot"
)

var (
	baselineID string
	currentID  string
	snapshotDir string
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare two log snapshots and show differences",
	RunE:  runDiff,
}

func init() {
	diffCmd.Flags().StringVar(&baselineID, "baseline", "", "ID of the baseline snapshot (required)")
	diffCmd.Flags().StringVar(&currentID, "current", "", "ID of the current snapshot (required)")
	diffCmd.Flags().StringVar(&snapshotDir, "dir", "snapshots", "Directory where snapshots are stored")
	_ = diffCmd.MarkFlagRequired("baseline")
	_ = diffCmd.MarkFlagRequired("current")
	RootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	baselinePath := fmt.Sprintf("%s/%s.json", snapshotDir, baselineID)
	currentPath := fmt.Sprintf("%s/%s.json", snapshotDir, currentID)

	baseline, err := snapshot.Load(baselinePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading baseline snapshot: %v\n", err)
		return err
	}

	current, err := snapshot.Load(currentPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading current snapshot: %v\n", err)
		return err
	}

	result := diff.Compare(baseline, current)

	fmt.Printf("Diff: %s (baseline) vs %s (current)\n", baselineID, currentID)
	fmt.Println(result.Summary())

	if len(result.Added) == 0 && len(result.Removed) == 0 && len(result.Changed) == 0 {
		fmt.Println("No differences found.")
	}

	return nil
}
