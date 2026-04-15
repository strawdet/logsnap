package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writePruneSnapshot(t *testing.T, dir, id string, createdAt time.Time) {
	t.Helper()
	s := &Snapshot{
		ID:        id,
		CreatedAt: createdAt,
		Entries:   []LogEntry{{Level: "info", Message: "test"}},
	}
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, id+".json"), data, 0644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}
}

func TestPrune_KeepLast(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	writePruneSnapshot(t, dir, "snap-1", now.Add(-3*time.Hour))
	writePruneSnapshot(t, dir, "snap-2", now.Add(-2*time.Hour))
	writePruneSnapshot(t, dir, "snap-3", now.Add(-1*time.Hour))

	result, err := Prune(dir, PruneOptions{KeepLast: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Kept) != 2 {
		t.Errorf("expected 2 kept, got %d", len(result.Kept))
	}
	if len(result.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(result.Removed))
	}
}

func TestPrune_OlderThan(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	writePruneSnapshot(t, dir, "snap-old", now.Add(-48*time.Hour))
	writePruneSnapshot(t, dir, "snap-new", now.Add(-1*time.Hour))

	cutoff := now.Add(-24 * time.Hour)
	result, err := Prune(dir, PruneOptions{OlderThan: cutoff})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 1 || result.Removed[0] != "snap-old" {
		t.Errorf("expected snap-old removed, got %v", result.Removed)
	}
}

func TestPrune_DryRun(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	writePruneSnapshot(t, dir, "snap-a", now.Add(-2*time.Hour))
	writePruneSnapshot(t, dir, "snap-b", now.Add(-1*time.Hour))

	result, err := Prune(dir, PruneOptions{KeepLast: 1, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 1 {
		t.Errorf("expected 1 in removed list, got %d", len(result.Removed))
	}
	// File should still exist after dry run.
	path := filepath.Join(dir, result.Removed[0]+".json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("dry run should not delete files")
	}
}

func TestPrune_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	result, err := Prune(dir, PruneOptions{KeepLast: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Removed) != 0 || len(result.Kept) != 0 {
		t.Error("expected empty result for empty dir")
	}
}
