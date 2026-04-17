package snapshot

import (
	"os"
	"testing"
)

func writeNoteSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	s := &Snapshot{ID: id, Entries: []LogEntry{}}
	if err := s.Save(dir); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
}

func TestAddNote_AndGet(t *testing.T) {
	dir := t.TempDir()
	writeNoteSnapshot(t, dir, "snap1")
	if err := AddNote(dir, "snap1", "initial note"); err != nil {
		t.Fatalf("AddNote: %v", err)
	}
	note, err := GetNote(dir, "snap1")
	if err != nil {
		t.Fatalf("GetNote: %v", err)
	}
	if note.Text != "initial note" {
		t.Errorf("expected 'initial note', got %q", note.Text)
	}
}

func TestAddNote_PreservesCreatedAt(t *testing.T) {
	dir := t.TempDir()
	writeNoteSnapshot(t, dir, "snap2")
	AddNote(dir, "snap2", "first")
	n1, _ := GetNote(dir, "snap2")
	AddNote(dir, "snap2", "updated")
	n2, _ := GetNote(dir, "snap2")
	if !n1.CreatedAt.Equal(n2.CreatedAt) {
		t.Error("CreatedAt should be preserved on update")
	}
	if n2.Text != "updated" {
		t.Errorf("expected 'updated', got %q", n2.Text)
	}
}

func TestAddNote_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := AddNote(dir, "missing", "text")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestGetNote_NoNote(t *testing.T) {
	dir := t.TempDir()
	writeNoteSnapshot(t, dir, "snap3")
	_, err := GetNote(dir, "snap3")
	if err == nil {
		t.Error("expected error when no note exists")
	}
}

func TestRemoveNote(t *testing.T) {
	dir := t.TempDir()
	writeNoteSnapshot(t, dir, "snap4")
	AddNote(dir, "snap4", "to remove")
	if err := RemoveNote(dir, "snap4"); err != nil {
		t.Fatalf("RemoveNote: %v", err)
	}
	if _, err := os.Stat(notePath(dir, "snap4")); !os.IsNotExist(err) {
		t.Error("note file should be removed")
	}
}

func TestRemoveNote_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := RemoveNote(dir, "ghost")
	if err == nil {
		t.Error("expected error removing non-existent note")
	}
}
