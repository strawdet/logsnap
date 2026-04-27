package snapshot

import (
	"testing"
	"time"
)

func writeVersionSnapshot(t *testing.T, dir string) *Snapshot {
	t.Helper()
	snap := &Snapshot{
		ID:        "ver-snap-001",
		Label:     "version-test",
		CreatedAt: time.Now().UTC(),
		Entries:   []LogEntry{{Level: "info", Message: "boot", Timestamp: time.Now().UTC()}},
	}
	if err := snap.Save(dir); err != nil {
		t.Fatalf("failed to write snapshot: %v", err)
	}
	return snap
}

func TestAddVersion_AndLoad(t *testing.T) {
	dir := t.TempDir()
	writeVersionSnapshot(t, dir)

	if err := AddVersion(dir, "ver-snap-001", "initial release"); err != nil {
		t.Fatalf("AddVersion failed: %v", err)
	}
	if err := AddVersion(dir, "ver-snap-001", "hotfix patch"); err != nil {
		t.Fatalf("AddVersion second call failed: %v", err)
	}

	index, err := LoadVersionIndex(dir, "ver-snap-001")
	if err != nil {
		t.Fatalf("LoadVersionIndex failed: %v", err)
	}
	if len(index.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(index.Entries))
	}
	if index.Entries[0].Version != 1 {
		t.Errorf("expected version 1, got %d", index.Entries[0].Version)
	}
	if index.Entries[1].Note != "hotfix patch" {
		t.Errorf("unexpected note: %s", index.Entries[1].Note)
	}
	if index.Entries[1].Label != "version-test" {
		t.Errorf("expected label 'version-test', got %s", index.Entries[1].Label)
	}
}

func TestAddVersion_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := AddVersion(dir, "nonexistent", "note")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestLoadVersionIndex_NoFile(t *testing.T) {
	dir := t.TempDir()
	index, err := LoadVersionIndex(dir, "missing-snap")
	if err != nil {
		t.Fatalf("expected no error for missing index, got: %v", err)
	}
	if len(index.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(index.Entries))
	}
}

func TestClearVersionHistory(t *testing.T) {
	dir := t.TempDir()
	writeVersionSnapshot(t, dir)

	_ = AddVersion(dir, "ver-snap-001", "v1")

	if err := ClearVersionHistory(dir, "ver-snap-001"); err != nil {
		t.Fatalf("ClearVersionHistory failed: %v", err)
	}

	index, err := LoadVersionIndex(dir, "ver-snap-001")
	if err != nil {
		t.Fatalf("unexpected error after clear: %v", err)
	}
	if len(index.Entries) != 0 {
		t.Errorf("expected empty index after clear, got %d entries", len(index.Entries))
	}
}

func TestClearVersionHistory_NoFile(t *testing.T) {
	dir := t.TempDir()
	if err := ClearVersionHistory(dir, "ghost"); err != nil {
		t.Errorf("expected no error clearing nonexistent history, got: %v", err)
	}
}
