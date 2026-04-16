package snapshot_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/petems/logsnap/internal/snapshot"
)

func writeRenameSnapshot(t *testing.T, dir, id, label string) {
	t.Helper()
	snap := snapshot.Snapshot{
		ID:        id,
		Label:     label,
		CreatedAt: time.Now(),
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, id+".json"), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestRenameSnapshot(t *testing.T) {
	dir := t.TempDir()
	writeRenameSnapshot(t, dir, "snap-001", "old-label")

	if err := snapshot.RenameSnapshot(dir, "snap-001", "new-label"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap, err := snapshot.Load(dir, "snap-001")
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if snap.Label != "new-label" {
		t.Errorf("expected label %q, got %q", "new-label", snap.Label)
	}
}

func TestRenameSnapshot_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := snapshot.RenameSnapshot(dir, "missing", "label")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestRenameSnapshot_PreservesEntries(t *testing.T) {
	dir := t.TempDir()
	writeRenameSnapshot(t, dir, "snap-002", "original")

	if err := snapshot.RenameSnapshot(dir, "snap-002", "updated"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snap, err := snapshot.Load(dir, "snap-002")
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if snap.ID != "snap-002" {
		t.Errorf("expected ID %q, got %q", "snap-002", snap.ID)
	}
}
