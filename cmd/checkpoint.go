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

	checkpointCmd := &cobra.Command{
		Use:   "checkpoint",
		Short: "Manage named checkpoints that reference snapshots",
	}

	setCmd := &cobra.Command{
		Use:   "set <name> <snapshot-id>",
		Short: "Create or update a checkpoint",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			desc, _ := cmd.Flags().GetString("desc")
			if err := snapshot.SetCheckpoint(dir, args[0], args[1], desc); err != nil {
				return fmt.Errorf("set checkpoint: %w", err)
			}
			fmt.Printf("Checkpoint %q set → %s\n", args[0], args[1])
			return nil
		},
	}
	setCmd.Flags().String("desc", "", "Optional description for the checkpoint")

	getCmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Resolve a checkpoint to its snapshot ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := snapshot.ResolveCheckpoint(dir, args[0])
			if err != nil {
				return err
			}
			fmt.Println(id)
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Delete a checkpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.RemoveCheckpoint(dir, args[0]); err != nil {
				return err
			}
			fmt.Printf("Checkpoint %q removed\n", args[0])
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all checkpoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			cps, err := snapshot.ListCheckpoints(dir)
			if err != nil {
				return err
			}
			if len(cps) == 0 {
				fmt.Println("No checkpoints found.")
				return nil
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tSNAPSHOT ID\tDESCRIPTION\tCREATED")
			for _, cp := range cps {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					cp.Name, cp.SnapshotID, cp.Description,
					cp.CreatedAt.Format("2006-01-02 15:04:05"))
			}
			return w.Flush()
		},
	}

	checkpointCmd.PersistentFlags().StringVar(&dir, "dir", ".logsnap", "Snapshot storage directory")
	checkpointCmd.AddCommand(setCmd, getCmd, removeCmd, listCmd)
	rootCmd.AddCommand(checkpointCmd)
}
