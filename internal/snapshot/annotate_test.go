package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func writeAnnotateSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	f := filepath.Join(dir, id+".json")
	if err := os.WriteFile(f, []byte(`{"id":"`+id+`"}`), 0644); err != nil {
		t.Fatalf("writeAnnotateSnapshot: %v", err)
	}
}

func TestAddAnnotation_AndGet(t *testing.T) {
	dir := t.TempDir()
	const id = "snap-abc123"
	writeAnnotateSnapshot(t, dir, id)

	if err := AddAnnotation(dir, id, "deployed to staging"); err != nil {
		t.Fatalf("AddAnnotation: %v", err)
	}

	a, err := GetAnnotation(dir, id)
	if err != nil {
		t.Fatalf("GetAnnotation: %v", err)
	}
	if a == nil {
		t.Fatal("expected annotation, got nil")
	}
	if a.Note != "deployed to staging" {
		t.Errorf("note = %q, want %q", a.Note, "deployed to staging")
	}
	if a.SnapshotID != id {
		t.Errorf("snapshot_id = %q, want %q", a.SnapshotID, id)
	}
}

func TestAddAnnotation_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := AddAnnotation(dir, "nonexistent", "some note")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestGetAnnotation_NoAnnotation(t *testing.T) {
	dir := t.TempDir()
	const id = "snap-xyz"
	writeAnnotateSnapshot(t, dir, id)

	a, err := GetAnnotation(dir, id)
	if err != nil {
		t.Fatalf("GetAnnotation: %v", err)
	}
	if a != nil {
		t.Errorf("expected nil annotation, got %+v", a)
	}
}

func TestRemoveAnnotation(t *testing.T) {
	dir := t.TempDir()
	const id = "snap-del"
	writeAnnotateSnapshot(t, dir, id)

	if err := AddAnnotation(dir, id, "to be removed"); err != nil {
		t.Fatalf("AddAnnotation: %v", err)
	}
	if err := RemoveAnnotation(dir, id); err != nil {
		t.Fatalf("RemoveAnnotation: %v", err)
	}

	a, err := GetAnnotation(dir, id)
	if err != nil {
		t.Fatalf("GetAnnotation after remove: %v", err)
	}
	if a != nil {
		t.Error("expected nil after removal")
	}
}

func TestRemoveAnnotation_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := RemoveAnnotation(dir, "no-such-snap")
	if err == nil {
		t.Fatal("expected error removing non-existent annotation")
	}
}
