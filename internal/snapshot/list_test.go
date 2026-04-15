package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

func TestListSnapshots_Empty(t *testing.T) {
	dir := t.TempDir()
	metas, err := snapshot.ListSnapshots(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("expected 0 snapshots, got %d", len(metas))
	}
}

func TestListSnapshots_NonExistentDir(t *testing.T) {
	metas, err := snapshot.ListSnapshots("/tmp/logsnap_does_not_exist_xyz")
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got: %v", err)
	}
	if metas != nil {
		t.Errorf("expected nil slice, got %v", metas)
	}
}

func TestListSnapshots_SortedNewestFirst(t *testing.T) {
	dir := t.TempDir()

	old := snapshot.New("old", sampleEntries())
	old.CreatedAt = time.Now().Add(-2 * time.Hour)
	if err := old.Save(dir); err != nil {
		t.Fatal(err)
	}

	recent := snapshot.New("recent", sampleEntries())
	recent.CreatedAt = time.Now()
	if err := recent.Save(dir); err != nil {
		t.Fatal(err)
	}

	metas, err := snapshot.ListSnapshots(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(metas))
	}
	if metas[0].Label != "recent" {
		t.Errorf("expected newest first, got label %q", metas[0].Label)
	}
}

func TestDeleteSnapshot(t *testing.T) {
	dir := t.TempDir()

	snap := snapshot.New("to-delete", sampleEntries())
	if err := snap.Save(dir); err != nil {
		t.Fatal(err)
	}

	if err := snapshot.DeleteSnapshot(dir, snap.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	metas, _ := snapshot.ListSnapshots(dir)
	if len(metas) != 0 {
		t.Errorf("expected snapshot to be deleted, still have %d", len(metas))
	}
}

func TestDeleteSnapshot_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := snapshot.DeleteSnapshot(dir, "nonexistent-id")
	if err == nil {
		t.Error("expected error for missing snapshot, got nil")
	}
}

func TestListSnapshots_SkipsInvalidFiles(t *testing.T) {
	dir := t.TempDir()
	// write a broken JSON file
	if err := os.WriteFile(dir+"/bad.json", []byte("not-json{"), 0644); err != nil {
		t.Fatal(err)
	}
	metas, err := snapshot.ListSnapshots(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(metas) != 0 {
		t.Errorf("expected invalid files to be skipped, got %d entries", len(metas))
	}
}
