package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	nsCmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage snapshot namespaces",
	}

	addNsCmd := &cobra.Command{
		Use:   "add <namespace> <snapshot-id>",
		Short: "Add a snapshot to a namespace",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := defaultSnapshotDir()
			if err := snapshot.AddToNamespace(dir, args[0], args[1]); err != nil {
				return fmt.Errorf("add to namespace: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added %s to namespace %q\n", args[1], args[0])
			return nil
		},
	}

	removeNsCmd := &cobra.Command{
		Use:   "remove <namespace> <snapshot-id>",
		Short: "Remove a snapshot from a namespace",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := defaultSnapshotDir()
			if err := snapshot.RemoveFromNamespace(dir, args[0], args[1]); err != nil {
				return fmt.Errorf("remove from namespace: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Removed %s from namespace %q\n", args[1], args[0])
			return nil
		},
	}

	listNsCmd := &cobra.Command{
		Use:   "list",
		Short: "List all namespaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := defaultSnapshotDir()
			names, err := snapshot.ListNamespaces(dir)
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No namespaces found.")
				return nil
			}
			for _, ns := range names {
				ids, _ := snapshot.GetNamespaceSnapshots(dir, ns)
				fmt.Fprintf(cmd.OutOrStdout(), "%s (%d): %s\n", ns, len(ids), strings.Join(ids, ", "))
			}
			return nil
		},
	}

	showNsCmd := &cobra.Command{
		Use:   "show <namespace>",
		Short: "Show snapshots in a namespace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := defaultSnapshotDir()
			ids, err := snapshot.GetNamespaceSnapshots(dir, args[0])
			if err != nil {
				return err
			}
			if len(ids) == 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Namespace %q is empty.\n", args[0])
				return nil
			}
			for _, id := range ids {
				fmt.Fprintln(cmd.OutOrStdout(), id)
			}
			return nil
		},
	}

	nsCmd.AddCommand(addNsCmd, removeNsCmd, listNsCmd, showNsCmd)
	rootCmd.AddCommand(nsCmd)
	_ = os.MkdirAll(defaultSnapshotDir(), 0755)
}
