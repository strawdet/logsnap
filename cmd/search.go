package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

var (
	searchTag      string
	searchSince    string
	searchUntil    string
	searchLabelKey string
	searchLabelVal string
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search snapshots by tag, label, or time range",
		RunE:  runSearch,
	}

	searchCmd.Flags().StringVar(&searchTag, "tag", "", "Filter by tag")
	searchCmd.Flags().StringVar(&searchSince, "since", "", "Filter snapshots created after this time (RFC3339)")
	searchCmd.Flags().StringVar(&searchUntil, "until", "", "Filter snapshots created before this time (RFC3339)")
	searchCmd.Flags().StringVar(&searchLabelKey, "label-key", "", "Filter by label key")
	searchCmd.Flags().StringVar(&searchLabelVal, "label-val", "", "Filter by label value (requires --label-key)")

	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	dir := snapshotDir()

	filter := snapshot.SearchFilter{
		Tag:      searchTag,
		LabelKey: searchLabelKey,
		LabelVal: searchLabelVal,
	}

	if searchSince != "" {
		t, err := time.Parse(time.RFC3339, searchSince)
		if err != nil {
			return fmt.Errorf("invalid --since format (use RFC3339): %w", err)
		}
		filter.Since = &t
	}

	if searchUntil != "" {
		t, err := time.Parse(time.RFC3339, searchUntil)
		if err != nil {
			return fmt.Errorf("invalid --until format (use RFC3339): %w", err)
		}
		filter.Until = &t
	}

	results, err := snapshot.Search(dir, filter)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		fmt.Println("No snapshots matched the given criteria.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tCREATED\tTAGS\tLABELS")
	for _, r := range results {
		s := r.Snapshot
		tags := joinStrings(s.Tags)
		labels := formatLabels(s.Labels)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", s.ID, s.CreatedAt.Format(time.RFC3339), tags, labels)
	}
	return w.Flush()
}

func joinStrings(ss []string) string {
	if len(ss) == 0 {
		return "-"
	}
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += ","
		}
		out += s
	}
	return out
}

func formatLabels(m map[string]string) string {
	if len(m) == 0 {
		return "-"
	}
	out := ""
	for k, v := range m {
		if out != "" {
			out += ","
		}
		out += k + "=" + v
	}
	return out
}
