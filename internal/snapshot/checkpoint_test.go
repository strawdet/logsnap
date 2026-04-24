package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCheckpointSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	f, err := os.Create(filepath.Join(dir, id+".json"))
	if err != nil {
		t.Fatalf("create snapshot file: %v", err)
	}
	f.WriteString(`{"id":"` + id + `","entries":[]}`)
	f.Close()
}

func TestSetCheckpoint_AndResolve(t *testing.T) {
	dir := t.TempDir()
	writeCheckpointSnapshot(t, dir, "snap-abc")

	if err := SetCheckpoint(dir, "v1", "snap-abc", "first release"); err != nil {
		t.Fatalf("SetCheckpoint: %v", err)
	}

	id, err := ResolveCheckpoint(dir, "v1")
	if err != nil {
		t.Fatalf("ResolveCheckpoint: %v", err)
	}
	if id != "snap-abc" {
		t.Errorf("expected snap-abc, got %s", id)
	}
}

func TestSetCheckpoint_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := SetCheckpoint(dir, "v1", "missing-snap", "")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestResolveCheckpoint_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := ResolveCheckpoint(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing checkpoint")
	}
}

func TestRemoveCheckpoint(t *testing.T) {
	dir := t.TempDir()
	writeCheckpointSnapshot(t, dir, "snap-xyz")

	if err := SetCheckpoint(dir, "release", "snap-xyz", ""); err != nil {
		t.Fatalf("SetCheckpoint: %v", err)
	}
	if err := RemoveCheckpoint(dir, "release"); err != nil {
		t.Fatalf("RemoveCheckpoint: %v", err)
	}
	_, err := ResolveCheckpoint(dir, "release")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestRemoveCheckpoint_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := RemoveCheckpoint(dir, "ghost")
	if err == nil {
		t.Fatal("expected error removing nonexistent checkpoint")
	}
}

func TestListCheckpoints(t *testing.T) {
	dir := t.TempDir()
	writeCheckpointSnapshot(t, dir, "snap-1")
	writeCheckpointSnapshot(t, dir, "snap-2")

	SetCheckpoint(dir, "alpha", "snap-1", "alpha build")
	SetCheckpoint(dir, "beta", "snap-2", "beta build")

	cps, err := ListCheckpoints(dir)
	if err != nil {
		t.Fatalf("ListCheckpoints: %v", err)
	}
	if len(cps) != 2 {
		t.Errorf("expected 2 checkpoints, got %d", len(cps))
	}
}

func TestLoadCheckpointIndex_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	idx, err := LoadCheckpointIndex(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx) != 0 {
		t.Errorf("expected empty index, got %d entries", len(idx))
	}
}
