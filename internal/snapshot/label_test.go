package snapshot

import (
	"os"
	"testing"
	"time"
)

func writeLabelSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	snap := &Snapshot{
		ID:        id,
		CreatedAt: time.Now(),
		Entries:   []LogEntry{{Level: "info", Message: "test"}},
	}
	if err := snap.Save(dir); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
}

func TestAddLabel_AndGet(t *testing.T) {
	dir := t.TempDir()
	writeLabelSnapshot(t, dir, "snap-001")

	if err := AddLabel(dir, "snap-001", "production"); err != nil {
		t.Fatalf("AddLabel: %v", err)
	}

	ids, err := GetSnapshotsByLabel(dir, "production")
	if err != nil {
		t.Fatalf("GetSnapshotsByLabel: %v", err)
	}
	if len(ids) != 1 || ids[0] != "snap-001" {
		t.Errorf("expected [snap-001], got %v", ids)
	}
}

func TestAddLabel_Deduplicates(t *testing.T) {
	dir := t.TempDir()
	writeLabelSnapshot(t, dir, "snap-001")

	_ = AddLabel(dir, "snap-001", "staging")
	_ = AddLabel(dir, "snap-001", "staging")

	ids, _ := GetSnapshotsByLabel(dir, "staging")
	if len(ids) != 1 {
		t.Errorf("expected 1 entry, got %d", len(ids))
	}
}

func TestAddLabel_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := AddLabel(dir, "nonexistent", "prod")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestRemoveLabel(t *testing.T) {
	dir := t.TempDir()
	writeLabelSnapshot(t, dir, "snap-001")
	_ = AddLabel(dir, "snap-001", "v1")

	if err := RemoveLabel(dir, "snap-001", "v1"); err != nil {
		t.Fatalf("RemoveLabel: %v", err)
	}

	ids, _ := GetSnapshotsByLabel(dir, "v1")
	if len(ids) != 0 {
		t.Errorf("expected empty, got %v", ids)
	}
}

func TestListLabelsForSnapshot(t *testing.T) {
	dir := t.TempDir()
	writeLabelSnapshot(t, dir, "snap-001")
	_ = AddLabel(dir, "snap-001", "prod")
	_ = AddLabel(dir, "snap-001", "v2")

	labels, err := ListLabelsForSnapshot(dir, "snap-001")
	if err != nil {
		t.Fatalf("ListLabelsForSnapshot: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d: %v", len(labels), labels)
	}
}

func TestLoadLabelIndex_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	idx, err := LoadLabelIndex(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx.Labels) != 0 {
		t.Errorf("expected empty index")
	}
}

func TestLoadLabelIndex_NonExistentDir(t *testing.T) {
	idx, err := LoadLabelIndex("/nonexistent/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx == nil {
		t.Error("expected non-nil index")
	}
}

func TestLoadLabelIndex_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(labelIndexPath(dir), []byte("not-json"), 0644)
	_, err := LoadLabelIndex(dir)
	if err == nil {
		t.Error("expected parse error")
	}
}
