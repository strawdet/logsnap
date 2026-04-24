package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func init() {
	var dir string

	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage snapshot workflows",
	}

	saveCmd := &cobra.Command{
		Use:   "save <name> <action[:param=val,...]>...",
		Short: "Save a new workflow",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			var steps []snapshot.WorkflowStep
			for _, raw := range args[1:] {
				step, err := parseStep(raw)
				if err != nil {
					return err
				}
				steps = append(steps, step)
			}
			wf := snapshot.Workflow{Name: name, Steps: steps}
			if err := snapshot.SaveWorkflow(dir, wf); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "workflow %q saved (%d steps)\n", name, len(steps))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved workflows",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := snapshot.ListWorkflows(dir)
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Fprintln(os.Stdout, "no workflows found")
				return nil
			}
			for _, n := range names {
				fmt.Fprintln(os.Stdout, n)
			}
			return nil
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete a workflow",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := snapshot.DeleteWorkflow(dir, args[0]); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "workflow %q deleted\n", args[0])
			return nil
		},
	}

	workflowCmd.PersistentFlags().StringVar(&dir, "dir", defaultSnapshotDir(), "snapshot directory")
	workflowCmd.AddCommand(saveCmd, listCmd, deleteCmd)
	RootCmd.AddCommand(workflowCmd)
}

// parseStep parses "action" or "action:key=val,key2=val2" into a WorkflowStep.
func parseStep(raw string) (snapshot.WorkflowStep, error) {
	parts := strings.SplitN(raw, ":", 2)
	step := snapshot.WorkflowStep{
		Name:   parts[0],
		Action: parts[0],
		Params: map[string]string{},
	}
	if len(parts) == 2 {
		for _, kv := range strings.Split(parts[1], ",") {
			pair := strings.SplitN(kv, "=", 2)
			if len(pair) != 2 {
				return step, fmt.Errorf("invalid param %q (expected key=value)", kv)
			}
			step.Params[pair[0]] = pair[1]
		}
	}
	return step, nil
}
