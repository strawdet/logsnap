package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"logsnap/internal/snapshot"
)

func writeBaselineSnap(t *testing.T, dir, id string) {
	t.Helper()
	snap := &snapshot.Snapshot{ID: id, Entries: []snapshot.LogEntry{}}
	data, _ := json.Marshal(snap)
	_ = os.MkdirAll(dir, 0755)
	_ = os.WriteFile(filepath.Join(dir, id+".json"), data, 0644)
}

func TestBaselineCommand_SetAndGet(t *testing.T) {
	dir := t.TempDir()
	writeBaselineSnap(t, dir, "snap-001")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	rootCmd.SetArgs([]string{"baseline", "set", "prod", "snap-001", "--dir", dir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("set baseline: %v", err)
	}

	buf.Reset()
	rootCmd.SetArgs([]string{"baseline", "get", "prod", "--dir", dir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("get baseline: %v", err)
	}
	if !strings.Contains(buf.String(), "snap-001") {
		t.Errorf("expected snap-001 in output, got: %s", buf.String())
	}
}

func TestBaselineCommand_List(t *testing.T) {
	dir := t.TempDir()
	writeBaselineSnap(t, dir, "snap-002")
	writeBaselineSnap(t, dir, "snap-003")

	_ = snapshot.SetBaseline(dir, "staging", "snap-002")
	_ = snapshot.SetBaseline(dir, "canary", "snap-003")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"baseline", "list", "--dir", dir})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list baselines: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "staging") || !strings.Contains(out, "canary") {
		t.Errorf("unexpected list output: %s", out)
	}
}
