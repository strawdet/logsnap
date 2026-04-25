package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	var snapshotDir string

	triggerCmd := &cobra.Command{
		Use:   "trigger",
		Short: "Manage auto-capture trigger rules",
	}

	// trigger save
	var (
		errorRate  float64
		minEntries int
		label      string
	)
	saveCmd := &cobra.Command{
		Use:   "save <name>",
		Short: "Save a trigger rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := snapshotDir
			if dir == "" {
				dir = defaultSnapshotDir()
			}
			cond := snapshot.TriggerCondition{
				ErrorRateThreshold: errorRate,
				MinEntries:         minEntries,
				Label:              label,
			}
			tr, err := snapshot.SaveTrigger(dir, args[0], cond)
			if err != nil {
				return fmt.Errorf("save trigger: %w", err)
			}
			fmt.Printf("Trigger %q saved (id=%s)\n", tr.Name, tr.ID)
			return nil
		},
	}
	saveCmd.Flags().Float64Var(&errorRate, "error-rate", 0, "Minimum error rate to fire (0.0–1.0)")
	saveCmd.Flags().IntVar(&minEntries, "min-entries", 0, "Minimum number of log entries")
	saveCmd.Flags().StringVar(&label, "label", "", "Label to attach to triggered snapshot")

	// trigger show
	showCmd := &cobra.Command{
		Use:   "show <name>",
		Short: "Show a trigger rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := snapshotDir
			if dir == "" {
				dir = defaultSnapshotDir()
			}
			tr, err := snapshot.LoadTrigger(dir, args[0])
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintf(w, "Name:\t%s\n", tr.Name)
			fmt.Fprintf(w, "ID:\t%s\n", tr.ID)
			fmt.Fprintf(w, "Created:\t%s\n", tr.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Fprintf(w, "Error Rate Threshold:\t%s\n", strconv.FormatFloat(tr.Condition.ErrorRateThreshold, 'f', 2, 64))
			fmt.Fprintf(w, "Min Entries:\t%d\n", tr.Condition.MinEntries)
			fmt.Fprintf(w, "Label:\t%s\n", tr.Condition.Label)
			return w.Flush()
		},
	}

	// trigger delete
	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a trigger rule",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := snapshotDir
			if dir == "" {
				dir = defaultSnapshotDir()
			}
			if err := snapshot.DeleteTrigger(dir, args[0]); err != nil {
				return err
			}
			fmt.Printf("Trigger %q deleted\n", args[0])
			return nil
		},
	}

	triggerCmd.PersistentFlags().StringVar(&snapshotDir, "dir", "", "Snapshot directory")
	triggerCmd.AddCommand(saveCmd, showCmd, deleteCmd)
	rootCmd.AddCommand(triggerCmd)
}
