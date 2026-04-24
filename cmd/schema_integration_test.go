package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"logsnap/internal/snapshot"
)

func writeSchemaSnapshot(t *testing.T, dir, id string, entries []snapshot.LogEntry) {
	t.Helper()
	snap := &snapshot.Snapshot{
		ID:        id,
		CreatedAt: time.Now(),
		Entries:   entries,
	}
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, id+".json"), data, 0644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}
}

func TestSchemaCommand_SaveAndShow(t *testing.T) {
	dir := t.TempDir()

	rootCmd.ResetFlags()
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)

	rootCmd.SetArgs([]string{
		"schema", "save", "myschema",
		"--dir", dir,
		"--required", "message,level",
		"--allowed-levels", "info,error",
		"--description", "test schema",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("schema save: %v", err)
	}

	s, err := snapshot.LoadSchema(dir, "myschema")
	if err != nil {
		t.Fatalf("LoadSchema: %v", err)
	}
	if s.Name != "myschema" {
		t.Errorf("expected name 'myschema', got %q", s.Name)
	}
	if len(s.AllowedLevels) != 2 {
		t.Errorf("expected 2 allowed levels, got %d", len(s.AllowedLevels))
	}
}

func TestSchemaCommand_Validate_Pass(t *testing.T) {
	dir := t.TempDir()
	writeSchemaSnapshot(t, dir, "snap-ok", []snapshot.LogEntry{
		{Message: "hello", Level: "info"},
	})

	s := &snapshot.Schema{
		ID: "v1", Name: "strict",
		RequiredFields: []string{"message", "level"},
		AllowedLevels:  []string{"info", "error"},
	}
	if err := snapshot.SaveSchema(dir, s); err != nil {
		t.Fatalf("SaveSchema: %v", err)
	}

	rootCmd.SetArgs([]string{"schema", "validate", "snap-ok", "strict", "--dir", dir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("validate should pass: %v", err)
	}
}

func TestSchemaCommand_Delete(t *testing.T) {
	dir := t.TempDir()
	s := &snapshot.Schema{ID: "d1", Name: "toDelete", RequiredFields: []string{"message"}}
	if err := snapshot.SaveSchema(dir, s); err != nil {
		t.Fatalf("SaveSchema: %v", err)
	}
	rootCmd.SetArgs([]string{"schema", "delete", "toDelete", "--dir", dir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := snapshot.LoadSchema(dir, "toDelete"); err == nil {
		t.Fatal("expected schema to be deleted")
	}
}
