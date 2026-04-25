package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

var pinNote string

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin <snapshot-id>",
		Short: "Pin a snapshot to prevent accidental deletion",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := cmd.Flag("dir").Value.String()
			if err := snapshot.PinSnapshot(dir, args[0], pinNote); err != nil {
				return fmt.Errorf("pin failed: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Pinned snapshot %s\n", args[0])
			return nil
		},
	}
	pinCmd.Flags().StringVarP(&pinNote, "note", "n", "", "Optional note for the pin")
	pinCmd.Flags().String("dir", "snapshots", "Directory where snapshots are stored")

	unpinCmd := &cobra.Command{
		Use:   "unpin <snapshot-id>",
		Short: "Unpin a previously pinned snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := cmd.Flag("dir").Value.String()
			if err := snapshot.UnpinSnapshot(dir, args[0]); err != nil {
				return fmt.Errorf("unpin failed: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Unpinned snapshot %s\n", args[0])
			return nil
		},
	}
	unpinCmd.Flags().String("dir", "snapshots", "Directory where snapshots are stored")

	pinsCmd := &cobra.Command{
		Use:   "pins",
		Short: "List all pinned snapshots",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := cmd.Flag("dir").Value.String()
			index, err := snapshot.LoadPinIndex(dir)
			if err != nil {
				return fmt.Errorf("could not load pins: %w", err)
			}
			if len(index) == 0 {
				fmt.Println("No pinned snapshots.")
				return nil
			}
			printPinIndex(index)
			return nil
		},
	}
	pinsCmd.Flags().String("dir", "snapshots", "Directory where snapshots are stored")

	rootCmd.AddCommand(pinCmd)
	rootCmd.AddCommand(unpinCmd)
	rootCmd.AddCommand(pinsCmd)
}

// printPinIndex writes all pinned snapshot entries to stdout, showing the
// optional note alongside the snapshot ID when one is present.
func printPinIndex(index map[string]string) {
	for id, note := range index {
		if note != "" {
			fmt.Printf("  %s  — %s\n", id, note)
		} else {
			fmt.Printf("  %s\n", id)
		}
	}
}
