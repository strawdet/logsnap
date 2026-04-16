package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeBaselineSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	snap := &Snapshot{ID: id, Entries: []LogEntry{}}
	data, _ := json.Marshal(snap)
	_ = os.MkdirAll(dir, 0755)
	_ = os.WriteFile(filepath.Join(dir, id+".json"), data, 0644)
}

func TestSetBaseline_AndResolve(t *testing.T) {
	dir := t.TempDir()
	writeBaselineSnapshot(t, dir, "snap-abc")

	if err := SetBaseline(dir, "production", "snap-abc"); err != nil {
		t.Fatalf("SetBaseline: %v", err)
	}
	id, err := ResolveBaseline(dir, "production")
	if err != nil {
		t.Fatalf("ResolveBaseline: %v", err)
	}
	if id != "snap-abc" {
		t.Errorf("expected snap-abc, got %s", id)
	}
}

func TestSetBaseline_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := SetBaseline(dir, "staging", "missing-id")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestResolveBaseline_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := ResolveBaseline(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing baseline")
	}
}

func TestRemoveBaseline(t *testing.T) {
	dir := t.TempDir()
	writeBaselineSnapshot(t, dir, "snap-xyz")

	_ = SetBaseline(dir, "dev", "snap-xyz")
	if err := RemoveBaseline(dir, "dev"); err != nil {
		t.Fatalf("RemoveBaseline: %v", err)
	}
	_, err := ResolveBaseline(dir, "dev")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestRemoveBaseline_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := RemoveBaseline(dir, "ghost")
	if err == nil {
		t.Fatal("expected error for missing baseline")
	}
}

func TestLoadBaselineIndex_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	idx, err := LoadBaselineIndex(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx) != 0 {
		t.Errorf("expected empty index, got %v", idx)
	}
}
