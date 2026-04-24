package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	var snapshotDir string

	timelineCmd := &cobra.Command{
		Use:   "timeline",
		Short: "Manage the event timeline for a snapshot",
	}

	addCmd := &cobra.Command{
		Use:   "add <snapshot-id> <event>",
		Short: "Add an event to a snapshot's timeline",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snapshotID := args[0]
			event := args[1]
			detail, _ := cmd.Flags().GetString("detail")
			if err := snapshot.AddTimelineEvent(snapshotDir, snapshotID, event, detail); err != nil {
				return fmt.Errorf("add timeline event: %w", err)
			}
			fmt.Printf("Event %q added to timeline for snapshot %q\n", event, snapshotID)
			return nil
		},
	}
	addCmd.Flags().String("detail", "", "Optional detail for the event")

	showCmd := &cobra.Command{
		Use:   "show <snapshot-id>",
		Short: "Show the timeline for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			tl, err := snapshot.GetTimeline(snapshotDir, args[0])
			if err != nil {
				return fmt.Errorf("get timeline: %w", err)
			}
			if len(tl.Entries) == 0 {
				fmt.Println("No timeline events recorded.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIMESTAMP\tEVENT\tDETAIL")
			for _, e := range tl.Entries {
				fmt.Fprintf(w, "%s\t%s\t%s\n", e.Timestamp.Format("2006-01-02 15:04:05"), e.Event, e.Detail)
			}
			return w.Flush()
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear <snapshot-id>",
		Short: "Clear the timeline for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.ClearTimeline(snapshotDir, args[0]); err != nil {
				return fmt.Errorf("clear timeline: %w", err)
			}
			fmt.Printf("Timeline cleared for snapshot %q\n", args[0])
			return nil
		},
	}

	for _, sub := range []*cobra.Command{addCmd, showCmd, clearCmd} {
		sub.Flags().StringVar(&snapshotDir, "dir", ".logsnap", "Directory where snapshots are stored")
		timelineCmd.AddCommand(sub)
	}

	rootCmd.AddCommand(timelineCmd)
}
