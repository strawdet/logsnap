package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	var outPath string

	archiveCmd := &cobra.Command{
		Use:   "archive [snapshot-ids...]",
		Short: "Archive snapshots into a zip file",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			if outPath == "" {
				return fmt.Errorf("--out is required")
			}
			if err := snapshot.ArchiveSnapshots(dir, args, outPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Archived %s → %s\n", strings.Join(args, ", "), outPath)
			return nil
		},
	}
	archiveCmd.Flags().StringVar(&outPath, "out", "", "Output zip file path")

	unarchiveCmd := &cobra.Command{
		Use:   "unarchive [archive-path]",
		Short: "Restore snapshots from a zip archive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			ids, err := snapshot.UnarchiveSnapshots(dir, args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Restored snapshots: %s\n", strings.Join(ids, ", "))
			return nil
		},
	}

	rootCmd.AddCommand(archiveCmd)
	rootCmd.AddCommand(unarchiveCmd)
}
