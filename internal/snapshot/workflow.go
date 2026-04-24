package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// WorkflowStep defines a single step in a snapshot workflow.
type WorkflowStep struct {
	Name   string            `json:"name"`
	Action string            `json:"action"` // e.g. "capture", "diff", "export"
	Params map[string]string `json:"params,omitempty"`
}

// Workflow represents a named sequence of steps to execute against snapshots.
type Workflow struct {
	Name      string         `json:"name"`
	Steps     []WorkflowStep `json:"steps"`
	CreatedAt time.Time      `json:"created_at"`
}

func workflowPath(dir, name string) string {
	return filepath.Join(dir, "workflows", name+".json")
}

// SaveWorkflow persists a workflow definition to disk.
func SaveWorkflow(dir string, wf Workflow) error {
	wf.CreatedAt = time.Now().UTC()
	p := workflowPath(dir, wf.Name)
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("create workflow: %w", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(wf)
}

// LoadWorkflow reads a workflow definition from disk.
func LoadWorkflow(dir, name string) (*Workflow, error) {
	p := workflowPath(dir, name)
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("workflow %q not found", name)
		}
		return nil, fmt.Errorf("open workflow: %w", err)
	}
	defer f.Close()
	var wf Workflow
	if err := json.NewDecoder(f).Decode(&wf); err != nil {
		return nil, fmt.Errorf("decode workflow: %w", err)
	}
	return &wf, nil
}

// ListWorkflows returns all saved workflow names.
func ListWorkflows(dir string) ([]string, error) {
	wdir := filepath.Join(dir, "workflows")
	entries, err := os.ReadDir(wdir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read workflows dir: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

// DeleteWorkflow removes a workflow definition from disk.
func DeleteWorkflow(dir, name string) error {
	p := workflowPath(dir, name)
	if err := os.Remove(p); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("workflow %q not found", name)
		}
		return fmt.Errorf("delete workflow: %w", err)
	}
	return nil
}
