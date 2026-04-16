package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeArchiveSnapshot(t *testing.T, dir, id, label string) *Snapshot {
	t.Helper()
	snap := &Snapshot{
		ID:        id,
		Label:     label,
		CreatedAt: time.Now(),
		Entries:   []LogEntry{{Level: "info", Message: "hello", Timestamp: time.Now()}},
	}
	if err := snap.Save(dir); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	return snap
}

func TestArchiveAndUnarchive(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	makeArchiveSnapshot(t, srcDir, "snap-001", "first")
	makeArchiveSnapshot(t, srcDir, "snap-002", "second")

	archivePath := filepath.Join(t.TempDir(), "export.zip")
	if err := ArchiveSnapshots(srcDir, []string{"snap-001", "snap-002"}, archivePath); err != nil {
		t.Fatalf("archive: %v", err)
	}

	if _, err := os.Stat(archivePath); err != nil {
		t.Fatalf("archive file missing: %v", err)
	}

	ids, err := UnarchiveSnapshots(dstDir, archivePath)
	if err != nil {
		t.Fatalf("unarchive: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d", len(ids))
	}

	for _, id := range []string{"snap-001", "snap-002"} {
		snap, err := Load(dstDir, id)
		if err != nil {
			t.Fatalf("load restored %s: %v", id, err)
		}
		if snap.ID != id {
			t.Errorf("expected id %s, got %s", id, snap.ID)
		}
	}
}

func TestArchiveSnapshots_NotFound(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(t.TempDir(), "out.zip")
	err := ArchiveSnapshots(dir, []string{"missing"}, out)
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestUnarchiveSnapshots_BadFile(t *testing.T) {
	bad := filepath.Join(t.TempDir(), "bad.zip")
	_ = os.WriteFile(bad, []byte("not a zip"), 0644)
	_, err := UnarchiveSnapshots(t.TempDir(), bad)
	if err == nil {
		t.Fatal("expected error for bad zip")
	}
}
