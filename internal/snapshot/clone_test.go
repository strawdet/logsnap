package snapshot

import (
	"os"
	"testing"
	"time"
)

func makeCloneSnapshot(t *testing.T, dir string) *Snapshot {
	t.Helper()
	snap := &Snapshot{
		ID:        "clone-src-001",
		Label:     "original",
		CreatedAt: time.Now().UTC(),
		Entries: []LogEntry{
			{Level: "INFO", Message: "startup"},
			{Level: "ERROR", Message: "disk full"},
		},
	}
	if err := snap.Save(dir); err != nil {
		t.Fatalf("setup: failed to save snapshot: %v", err)
	}
	return snap
}

func TestCloneSnapshot(t *testing.T) {
	dir := t.TempDir()
	src := makeCloneSnapshot(t, dir)

	cloned, err := CloneSnapshot(dir, src.ID, "my-clone")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cloned.ID == src.ID {
		t.Error("cloned snapshot should have a different ID")
	}
	if cloned.Label != "my-clone" {
		t.Errorf("expected label 'my-clone', got %q", cloned.Label)
	}
	if len(cloned.Entries) != len(src.Entries) {
		t.Errorf("expected %d entries, got %d", len(src.Entries), len(cloned.Entries))
	}

	loaded, err := Load(dir, cloned.ID)
	if err != nil {
		t.Fatalf("failed to load cloned snapshot: %v", err)
	}
	if loaded.Label != "my-clone" {
		t.Errorf("persisted label mismatch: %q", loaded.Label)
	}
}

func TestCloneSnapshot_DefaultLabel(t *testing.T) {
	dir := t.TempDir()
	src := makeCloneSnapshot(t, dir)

	cloned, err := CloneSnapshot(dir, src.ID, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := src.Label + " (clone)"
	if cloned.Label != expected {
		t.Errorf("expected %q, got %q", expected, cloned.Label)
	}
}

func TestCloneSnapshot_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := CloneSnapshot(dir, "nonexistent", "")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestCloneSnapshot_InvalidDir(t *testing.T) {
	dir := t.TempDir()
	src := makeCloneSnapshot(t, dir)

	badDir := dir + "/no/perms"
	os.MkdirAll(badDir, 0000)
	defer os.Chmod(badDir, 0755)

	// Should still work since source is in dir, not badDir
	_, err := CloneSnapshot(dir, src.ID, "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
