package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	var snapshotDir string

	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage snapshot groups",
	}

	addCmd := &cobra.Command{
		Use:   "add <group> <snapshot-id>",
		Short: "Add a snapshot to a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.AddToGroup(snapshotDir, args[0], args[1]); err != nil {
				return fmt.Errorf("add to group: %w", err)
			}
			fmt.Printf("Added %s to group %q\n", args[1], args[0])
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <group> <snapshot-id>",
		Short: "Remove a snapshot from a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.RemoveFromGroup(snapshotDir, args[0], args[1]); err != nil {
				return fmt.Errorf("remove from group: %w", err)
			}
			fmt.Printf("Removed %s from group %q\n", args[1], args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := snapshot.ListGroups(snapshotDir)
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Println("No groups found.")
				return nil
			}
			for _, name := range names {
				g, _ := snapshot.GetGroup(snapshotDir, name)
				fmt.Printf("%s (%d snapshots)\n", name, len(g.Snapshots))
			}
			return nil
		},
	}

	showCmd := &cobra.Command{
		Use:   "show <group>",
		Short: "Show snapshots in a group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, err := snapshot.GetGroup(snapshotDir, args[0])
			if err != nil {
				return err
			}
			if len(g.Snapshots) == 0 {
				fmt.Printf("Group %q is empty.\n", args[0])
				return nil
			}
			fmt.Printf("Group: %s\n", g.Name)
			fmt.Printf("Snapshots:\n  %s\n", strings.Join(g.Snapshots, "\n  "))
			return nil
		},
	}

	defaultDir := os.Getenv("LOGSNAP_DIR")
	if defaultDir == "" {
		defaultDir = ".logsnap"
	}
	for _, sub := range []*cobra.Command{addCmd, removeCmd, listCmd, showCmd} {
		sub.Flags().StringVarP(&snapshotDir, "dir", "d", defaultDir, "Snapshot storage directory")
		groupCmd.AddCommand(sub)
	}

	RootCmd.AddCommand(groupCmd)
}
