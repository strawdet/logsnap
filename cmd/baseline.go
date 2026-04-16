package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	var dir string

	baselineCmd := &cobra.Command{
		Use:   "baseline",
		Short: "Manage named baselines pointing to snapshots",
	}

	setCmd := &cobra.Command{
		Use:   "set <name> <snapshot-id>",
		Short: "Assign a snapshot to a named baseline",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.SetBaseline(dir, args[0], args[1]); err != nil {
				return fmt.Errorf("set baseline: %w", err)
			}
			fmt.Printf("Baseline %q -> %s\n", args[0], args[1])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Resolve a named baseline to its snapshot ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := snapshot.ResolveBaseline(dir, args[0])
			if err != nil {
				return fmt.Errorf("resolve baseline: %w", err)
			}
			fmt.Println(id)
			return nil
		},
	}

	rmCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a named baseline",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.RemoveBaseline(dir, args[0]); err != nil {
				return fmt.Errorf("remove baseline: %w", err)
			}
			fmt.Printf("Removed baseline %q\n", args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all named baselines",
		RunE: func(cmd *cobra.Command, args []string) error {
			idx, err := snapshot.LoadBaselineIndex(dir)
			if err != nil {
				return fmt.Errorf("load baseline index: %w", err)
			}
			if len(idx) == 0 {
				fmt.Println("No baselines defined.")
				return nil
			}
			for name, id := range idx {
				fmt.Printf("%-20s %s\n", name, id)
			}
			return nil
		},
	}

	defaultDir := os.Getenv("LOGSNAP_DIR")
	if defaultDir == "" {
		defaultDir = ".logsnap"
	}
	for _, sub := range []*cobra.Command{setCmd, getCmd, rmCmd, listCmd} {
		sub.Flags().StringVar(&dir, "dir", defaultDir, "snapshot storage directory")
		baselineCmd.AddCommand(sub)
	}

	rootCmd.AddCommand(baselineCmd)
}
