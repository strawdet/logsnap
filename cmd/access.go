package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	var dir string

	accessCmd := &cobra.Command{
		Use:   "access",
		Short: "Manage access logs for snapshots",
	}

	recordCmd := &cobra.Command{
		Use:   "record <snapshot-id>",
		Short: "Record an access event for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			action, _ := cmd.Flags().GetString("action")
			actor, _ := cmd.Flags().GetString("actor")
			note, _ := cmd.Flags().GetString("note")
			if action == "" || actor == "" {
				return fmt.Errorf("--action and --actor are required")
			}
			if err := snapshot.RecordAccess(dir, args[0], action, actor, note); err != nil {
				return fmt.Errorf("record access: %w", err)
			}
			fmt.Printf("Access event recorded for snapshot %s\n", args[0])
			return nil
		},
	}
	recordCmd.Flags().String("action", "", "Action performed (e.g. read, export, diff)")
	recordCmd.Flags().String("actor", "", "User or system performing the action")
	recordCmd.Flags().String("note", "", "Optional note")

	showCmd := &cobra.Command{
		Use:   "show <snapshot-id>",
		Short: "Show access log for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := snapshot.GetAccessLog(dir, args[0])
			if err != nil {
				return fmt.Errorf("get access log: %w", err)
			}
			if len(log.Events) == 0 {
				fmt.Println("No access events recorded.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIME\tACTION\tACTOR\tNOTE")
			for _, e := range log.Events {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.At.Format("2006-01-02 15:04:05"), e.Action, e.Actor, e.Note)
			}
			return w.Flush()
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear <snapshot-id>",
		Short: "Clear access log for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.ClearAccessLog(dir, args[0]); err != nil {
				return fmt.Errorf("clear access log: %w", err)
			}
			fmt.Printf("Access log cleared for snapshot %s\n", args[0])
			return nil
		},
	}

	accessCmd.PersistentFlags().StringVar(&dir, "dir", ".logsnap", "Snapshot storage directory")
	accessCmd.AddCommand(recordCmd, showCmd, clearCmd)
	RootCmd.AddCommand(accessCmd)
}
