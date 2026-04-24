package snapshot

import (
	"os"
	"testing"
)

func writeNamespaceSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	snap := &Snapshot{
		ID:      id,
		Entries: []LogEntry{{Level: "info", Message: "ns test"}},
	}
	if err := snap.Save(dir); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
}

func TestAddToNamespace_AndGet(t *testing.T) {
	dir := t.TempDir()
	writeNamespaceSnapshot(t, dir, "snap-ns-1")
	writeNamespaceSnapshot(t, dir, "snap-ns-2")

	if err := AddToNamespace(dir, "production", "snap-ns-1"); err != nil {
		t.Fatalf("AddToNamespace: %v", err)
	}
	if err := AddToNamespace(dir, "production", "snap-ns-2"); err != nil {
		t.Fatalf("AddToNamespace: %v", err)
	}

	ids, err := GetNamespaceSnapshots(dir, "production")
	if err != nil {
		t.Fatalf("GetNamespaceSnapshots: %v", err)
	}
	if len(ids) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(ids))
	}
}

func TestAddToNamespace_Deduplicates(t *testing.T) {
	dir := t.TempDir()
	writeNamespaceSnapshot(t, dir, "snap-dedup")

	_ = AddToNamespace(dir, "staging", "snap-dedup")
	_ = AddToNamespace(dir, "staging", "snap-dedup")

	ids, _ := GetNamespaceSnapshots(dir, "staging")
	if len(ids) != 1 {
		t.Errorf("expected 1, got %d", len(ids))
	}
}

func TestAddToNamespace_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := AddToNamespace(dir, "prod", "nonexistent")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestRemoveFromNamespace(t *testing.T) {
	dir := t.TempDir()
	writeNamespaceSnapshot(t, dir, "snap-rem")
	_ = AddToNamespace(dir, "dev", "snap-rem")

	if err := RemoveFromNamespace(dir, "dev", "snap-rem"); err != nil {
		t.Fatalf("RemoveFromNamespace: %v", err)
	}
	ids, _ := GetNamespaceSnapshots(dir, "dev")
	if len(ids) != 0 {
		t.Errorf("expected 0 snapshots after removal, got %d", len(ids))
	}
}

func TestListNamespaces(t *testing.T) {
	dir := t.TempDir()
	writeNamespaceSnapshot(t, dir, "snap-a")
	writeNamespaceSnapshot(t, dir, "snap-b")
	_ = AddToNamespace(dir, "alpha", "snap-a")
	_ = AddToNamespace(dir, "beta", "snap-b")

	names, err := ListNamespaces(dir)
	if err != nil {
		t.Fatalf("ListNamespaces: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 namespaces, got %d", len(names))
	}
	if names[0] != "alpha" || names[1] != "beta" {
		t.Errorf("unexpected order: %v", names)
	}
}

func TestLoadNamespaceIndex_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	idx, err := LoadNamespaceIndex(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx.Namespaces) != 0 {
		t.Errorf("expected empty index")
	}
	_ = os.Remove(dir)
}
