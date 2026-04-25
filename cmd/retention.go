package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	retentionCmd := &cobra.Command{
		Use:   "retention",
		Short: "Manage snapshot retention policies",
	}

	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Set the retention policy for the snapshot directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			maxCount, _ := cmd.Flags().GetInt("max-count")
			maxAge, _ := cmd.Flags().GetInt("max-age-days")
			protectPins, _ := cmd.Flags().GetBool("protect-pins")

			policy := snapshot.RetentionPolicy{
				MaxCount:    maxCount,
				MaxAgeDays:  maxAge,
				ProtectPins: protectPins,
			}
			if err := snapshot.SetRetentionPolicy(dir, policy); err != nil {
				return err
			}
			fmt.Println("Retention policy saved.")
			return nil
		},
	}
	setCmd.Flags().String("dir", defaultSnapshotDir(), "snapshot directory")
	setCmd.Flags().Int("max-count", 0, "maximum number of snapshots to keep (0 = unlimited)")
	setCmd.Flags().Int("max-age-days", 0, "delete snapshots older than N days (0 = disabled)")
	setCmd.Flags().Bool("protect-pins", true, "skip pinned snapshots during cleanup")

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show the current retention policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			policy, err := snapshot.GetRetentionPolicy(dir)
			if err != nil {
				return err
			}
			fmt.Printf("Max Count:    %d\n", policy.MaxCount)
			fmt.Printf("Max Age Days: %d\n", policy.MaxAgeDays)
			fmt.Printf("Protect Pins: %v\n", policy.ProtectPins)
			return nil
		},
	}
	showCmd.Flags().String("dir", defaultSnapshotDir(), "snapshot directory")

	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply the retention policy, deleting qualifying snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			deleted, err := snapshot.ApplyRetentionPolicy(dir, dryRun)
			if err != nil {
				return err
			}
			if len(deleted) == 0 {
				fmt.Println("No snapshots matched the retention policy.")
				return nil
			}
			verb := "Deleted"
			if dryRun {
				verb = "Would delete"
			}
			for _, id := range deleted {
				fmt.Fprintf(os.Stdout, "%s: %s\n", verb, id)
			}
			return nil
		},
	}
	applyCmd.Flags().String("dir", defaultSnapshotDir(), "snapshot directory")
	applyCmd.Flags().Bool("dry-run", false, "preview deletions without removing files")

	retentionCmd.AddCommand(setCmd, showCmd, applyCmd)
	rootCmd.AddCommand(retentionCmd)
}
