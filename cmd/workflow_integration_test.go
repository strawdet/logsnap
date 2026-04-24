package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func rootCmdWithWorkflow(dir string) *cobra.Command {
	root := &cobra.Command{Use: "logsnap"}

	wf := &cobra.Command{Use: "workflow"}

	save := &cobra.Command{
		Use:  "save",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var steps []snapshot.WorkflowStep
			for _, raw := range args[1:] {
				s, err := parseStep(raw)
				if err != nil {
					return err
				}
				steps = append(steps, s)
			}
			return snapshot.SaveWorkflow(dir, snapshot.Workflow{Name: args[0], Steps: steps})
		},
	}

	list := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := snapshot.ListWorkflows(dir)
			if err != nil {
				return err
			}
			for _, n := range names {
				cmd.Println(n)
			}
			return nil
		},
	}

	wf.AddCommand(save, list)
	root.AddCommand(wf)
	return root
}

func TestWorkflowCommand_SaveAndList(t *testing.T) {
	dir := t.TempDir()
	root := rootCmdWithWorkflow(dir)

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)

	root.SetArgs([]string{"workflow", "save", "deploy", "capture:label=prod", "diff"})
	if err := root.Execute(); err != nil {
		t.Fatalf("save: %v", err)
	}

	buf.Reset()
	root.SetArgs([]string{"workflow", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "deploy") {
		t.Errorf("expected 'deploy' in output, got: %q", out)
	}
}

func TestWorkflowCommand_ParseStep(t *testing.T) {
	step, err := parseStep("capture:label=prod,env=staging")
	if err != nil {
		t.Fatalf("parseStep: %v", err)
	}
	if step.Action != "capture" {
		t.Errorf("action: got %q", step.Action)
	}
	if step.Params["label"] != "prod" {
		t.Errorf("label param: got %q", step.Params["label"])
	}
	if step.Params["env"] != "staging" {
		t.Errorf("env param: got %q", step.Params["env"])
	}
}

func TestWorkflowCommand_SaveAndLoad_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	wf := snapshot.Workflow{
		Name: "ci",
		Steps: []snapshot.WorkflowStep{
			{Name: "capture", Action: "capture"},
			{Name: "export", Action: "export", Params: map[string]string{"format": "csv"}},
		},
	}
	if err := snapshot.SaveWorkflow(dir, wf); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := snapshot.LoadWorkflow(dir, "ci")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded.Steps[1].Params["format"] != "csv" {
		t.Errorf("format param: got %q", loaded.Steps[1].Params["format"])
	}
	_ = os.RemoveAll(dir)
}
