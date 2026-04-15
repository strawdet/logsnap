package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/yourorg/logsnap/internal/diff"
	"github.com/yourorg/logsnap/internal/snapshot"
)

var (
	compareName   string
	compareDir    string
	compareList   bool
)

func init() {
	cmd := &cobra.Command{
		Use:   "compare <base-id> <target-id>",
		Short: "Compare two snapshots and optionally save the result",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runCompare,
	}

	cmd.Flags().StringVarP(&compareName, "name", "n", "", "Save comparison under this name")
	cmd.Flags().StringVar(&compareDir, "compare-dir", ".logsnap/compares", "Directory to store named comparisons")
	cmd.Flags().BoolVarP(&compareList, "list", "l", false, "List saved comparisons")

	rootCmd.AddCommand(cmd)
}

func runCompare(cmd *cobra.Command, args []string) error {
	if compareList {
		return listCompareResults()
	}

	if len(args) < 2 {
		return fmt.Errorf("requires <base-id> and <target-id>")
	}

	baseID, targetID := args[0], args[1]
	dir := ".logsnap/snapshots"

	base, err := snapshot.Load(filepath.Join(dir, baseID+".json"))
	if err != nil {
		return fmt.Errorf("load base snapshot: %w", err)
	}

	target, err := snapshot.Load(filepath.Join(dir, targetID+".json"))
	if err != nil {
		return fmt.Errorf("load target snapshot: %w", err)
	}

	results := diff.Compare(base, target)
	fmt.Println(results.Summary())

	if compareName != "" {
		_, err := snapshot.SaveCompareResult(compareDir, compareName, baseID, targetID)
		if err != nil {
			return fmt.Errorf("save compare result: %w", err)
		}
		fmt.Printf("Comparison saved as %q\n", compareName)
	}

	return nil
}

func listCompareResults() error {
	names, err := snapshot.ListCompareResults(compareDir)
	if err != nil {
		return fmt.Errorf("list comparisons: %w", err)
	}

	if len(names) == 0 {
		fmt.Println("No saved comparisons found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tBASE\tTARGET")
	for _, name := range names {
		r, err := snapshot.LoadCompareResult(compareDir, name)
		if err != nil {
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Name, r.BaseID, r.TargetID)
	}
	return w.Flush()
}
