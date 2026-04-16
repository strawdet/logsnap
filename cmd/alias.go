package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"logsnap/internal/snapshot"
)

func init() {
	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage snapshot aliases",
	}

	setCmd := &cobra.Command{
		Use:   "set <alias> <snapshot-id>",
		Short: "Assign an alias to a snapshot",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			if err := snapshot.SetAlias(dir, args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Alias %q -> %s set.\n", args[0], args[1])
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <alias>",
		Short: "Resolve an alias to a snapshot ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			id, err := snapshot.ResolveAlias(dir, args[0])
			if err != nil {
				return err
			}
			fmt.Println(id)
			return nil
		},
	}

	rmCmd := &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			if err := snapshot.RemoveAlias(dir, args[0]); err != nil {
				return err
			}
			fmt.Printf("Alias %q removed.\n", args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, _ := cmd.Flags().GetString("dir")
			idx, err := snapshot.ListAliases(dir)
			if err != nil {
				return err
			}
			if len(idx) == 0 {
				fmt.Println("No aliases defined.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ALIAS\tSNAPSHOT ID")
			for alias, id := range idx {
				fmt.Fprintf(w, "%s\t%s\n", alias, id)
			}
			w.Flush()
			return nil
		},
	}

	for _, sub := range []*cobra.Command{setCmd, getCmd, rmCmd, listCmd} {
		sub.Flags().String("dir", "snapshots", "Snapshot storage directory")
		aliasCmd.AddCommand(sub)
	}
	RootCmd.AddCommand(aliasCmd)
}
