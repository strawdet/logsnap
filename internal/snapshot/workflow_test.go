package snapshot

import (
	"os"
	"testing"
)

func TestSaveAndLoadWorkflow(t *testing.T) {
	dir := t.TempDir()
	wf := Workflow{
		Name: "nightly",
		Steps: []WorkflowStep{
			{Name: "capture", Action: "capture", Params: map[string]string{"label": "nightly"}},
			{Name: "diff", Action: "diff"},
		},
	}
	if err := SaveWorkflow(dir, wf); err != nil {
		t.Fatalf("SaveWorkflow: %v", err)
	}
	loaded, err := LoadWorkflow(dir, "nightly")
	if err != nil {
		t.Fatalf("LoadWorkflow: %v", err)
	}
	if loaded.Name != wf.Name {
		t.Errorf("name: got %q, want %q", loaded.Name, wf.Name)
	}
	if len(loaded.Steps) != 2 {
		t.Errorf("steps: got %d, want 2", len(loaded.Steps))
	}
	if loaded.Steps[0].Params["label"] != "nightly" {
		t.Errorf("param label: got %q", loaded.Steps[0].Params["label"])
	}
}

func TestLoadWorkflow_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadWorkflow(dir, "missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListWorkflows(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"alpha", "beta", "gamma"} {
		if err := SaveWorkflow(dir, Workflow{Name: name}); err != nil {
			t.Fatalf("SaveWorkflow %q: %v", name, err)
		}
	}
	names, err := ListWorkflows(dir)
	if err != nil {
		t.Fatalf("ListWorkflows: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("got %d workflows, want 3", len(names))
	}
}

func TestListWorkflows_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	names, err := ListWorkflows(dir)
	if err != nil {
		t.Fatalf("ListWorkflows: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty, got %v", names)
	}
}

func TestDeleteWorkflow(t *testing.T) {
	dir := t.TempDir()
	wf := Workflow{Name: "temp"}
	_ = SaveWorkflow(dir, wf)
	if err := DeleteWorkflow(dir, "temp"); err != nil {
		t.Fatalf("DeleteWorkflow: %v", err)
	}
	_, err := LoadWorkflow(dir, "temp")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestDeleteWorkflow_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := DeleteWorkflow(dir, "ghost")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestListWorkflows_NonExistentDir(t *testing.T) {
	dir := t.TempDir()
	_ = os.RemoveAll(dir)
	names, err := ListWorkflows(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if names != nil {
		t.Errorf("expected nil, got %v", names)
	}
}
