package snapshot

import (
	"path/filepath"
	"testing"
	"time"
)

func makeMergeSnapshot(t *testing.T, dir, label string, entries []LogEntry) *Snapshot {
	t.Helper()
	snap := &Snapshot{
		ID:        generateID(),
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Entries:   entries,
	}
	if err := snap.Save(filepath.Join(dir, snap.ID+".json")); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	return snap
}

func TestMergeSnapshots_NoConflicts(t *testing.T) {
	dir := t.TempDir()
	base := makeMergeSnapshot(t, dir, "base", []LogEntry{
		{Level: "info", Message: "startup"},
	})
	other := makeMergeSnapshot(t, dir, "other", []LogEntry{
		{Level: "warn", Message: "disk low"},
	})

	res, err := MergeSnapshots(dir, base.ID, other.ID, "merged")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
	if len(res.Merged.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Merged.Entries))
	}
}

func TestMergeSnapshots_WithConflicts(t *testing.T) {
	dir := t.TempDir()
	base := makeMergeSnapshot(t, dir, "base", []LogEntry{
		{Level: "info", Message: "startup"},
		{Level: "error", Message: "crash"},
	})
	other := makeMergeSnapshot(t, dir, "other", []LogEntry{
		{Level: "warn", Message: "crash"},
		{Level: "info", Message: "shutdown"},
	})

	res, err := MergeSnapshots(dir, base.ID, other.ID, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "crash" {
		t.Errorf("expected conflict on 'crash', got %v", res.Conflicts)
	}
	if len(res.Merged.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.Merged.Entries))
	}
}

func TestMergeSnapshots_BaseNotFound(t *testing.T) {
	dir := t.TempDir()
	other := makeMergeSnapshot(t, dir, "other", []LogEntry{})
	_, err := MergeSnapshots(dir, "nonexistent", other.ID, "")
	if err == nil {
		t.Error("expected error for missing base snapshot")
	}
}
