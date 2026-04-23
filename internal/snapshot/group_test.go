package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func writeGroupSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	s := &Snapshot{ID: id, Label: "test"}
	data, _ := marshalSnapshot(s)
	if err := os.WriteFile(filepath.Join(dir, id+".json"), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestAddToGroup_AndGet(t *testing.T) {
	dir := t.TempDir()
	writeGroupSnapshot(t, dir, "snap-001")
	writeGroupSnapshot(t, dir, "snap-002")

	if err := AddToGroup(dir, "release-1.0", "snap-001"); err != nil {
		t.Fatalf("AddToGroup: %v", err)
	}
	if err := AddToGroup(dir, "release-1.0", "snap-002"); err != nil {
		t.Fatalf("AddToGroup: %v", err)
	}

	g, err := GetGroup(dir, "release-1.0")
	if err != nil {
		t.Fatalf("GetGroup: %v", err)
	}
	if len(g.Snapshots) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(g.Snapshots))
	}
}

func TestAddToGroup_Deduplicates(t *testing.T) {
	dir := t.TempDir()
	writeGroupSnapshot(t, dir, "snap-001")

	_ = AddToGroup(dir, "grp", "snap-001")
	_ = AddToGroup(dir, "grp", "snap-001")

	g, err := GetGroup(dir, "grp")
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Snapshots) != 1 {
		t.Errorf("expected 1 snapshot, got %d", len(g.Snapshots))
	}
}

func TestAddToGroup_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := AddToGroup(dir, "grp", "missing-snap")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestRemoveFromGroup(t *testing.T) {
	dir := t.TempDir()
	writeGroupSnapshot(t, dir, "snap-001")
	writeGroupSnapshot(t, dir, "snap-002")

	_ = AddToGroup(dir, "grp", "snap-001")
	_ = AddToGroup(dir, "grp", "snap-002")

	if err := RemoveFromGroup(dir, "grp", "snap-001"); err != nil {
		t.Fatalf("RemoveFromGroup: %v", err)
	}

	g, _ := GetGroup(dir, "grp")
	if len(g.Snapshots) != 1 || g.Snapshots[0] != "snap-002" {
		t.Errorf("unexpected snapshots after remove: %v", g.Snapshots)
	}
}

func TestListGroups(t *testing.T) {
	dir := t.TempDir()
	writeGroupSnapshot(t, dir, "snap-001")

	_ = AddToGroup(dir, "alpha", "snap-001")
	_ = AddToGroup(dir, "beta", "snap-001")

	names, err := ListGroups(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 groups, got %d", len(names))
	}
}

func TestGetGroup_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := GetGroup(dir, "nonexistent")
	if err == nil {
		t.Error("expected error for missing group")
	}
}
