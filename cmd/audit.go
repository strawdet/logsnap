package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Manage audit logs for snapshots",
	}

	showCmd := &cobra.Command{
		Use:   "show <snapshot-id>",
		Short: "Show the audit log for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			log, err := snapshot.GetAuditLog(dir, args[0])
			if err != nil {
				return fmt.Errorf("load audit log: %w", err)
			}
			if len(log.Events) == 0 {
				fmt.Println("No audit events recorded.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIMESTAMP\tACTION\tDETAIL")
			for _, e := range log.Events {
				fmt.Fprintf(w, "%s\t%s\t%s\n", e.Timestamp.Format("2006-01-02 15:04:05"), e.Action, e.Detail)
			}
			return w.Flush()
		},
	}
	showCmd.Flags().String("dir", ".logsnap", "Snapshot storage directory")

	recordCmd := &cobra.Command{
		Use:   "record <snapshot-id> <action>",
		Short: "Manually record an audit event for a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			detail, _ := cmd.Flags().GetString("detail")
			if err := snapshot.RecordAuditEvent(dir, args[0], args[1], detail); err != nil {
				return fmt.Errorf("record audit event: %w", err)
			}
			fmt.Printf("Audit event '%s' recorded for snapshot %s.\n", args[1], args[0])
			return nil
		},
	}
	recordCmd.Flags().String("dir", ".logsnap", "Snapshot storage directory")
	recordCmd.Flags().String("detail", "", "Optional detail for the event")

	clearCmd := &cobra.Command{
		Use:   "clear <snapshot-id>",
		Short: "Clear the audit log for a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			if err := snapshot.ClearAuditLog(dir, args[0]); err != nil {
				return fmt.Errorf("clear audit log: %w", err)
			}
			fmt.Printf("Audit log cleared for snapshot %s.\n", args[0])
			return nil
		},
	}
	clearCmd.Flags().String("dir", ".logsnap", "Snapshot storage directory")

	auditCmd.AddCommand(showCmd, recordCmd, clearCmd)
	rootCmd.AddCommand(auditCmd)
}
