package snapshot

import (
	"testing"
)

func writeBookmarkSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	s := &Snapshot{ID: id, Label: "bm-" + id, Entries: nil}
	if err := s.Save(dir); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
}

func TestAddBookmark_AndResolve(t *testing.T) {
	dir := t.TempDir()
	writeBookmarkSnapshot(t, dir, "snap1")

	if err := AddBookmark(dir, "mymark", "snap1"); err != nil {
		t.Fatalf("AddBookmark: %v", err)
	}
	id, err := ResolveBookmark(dir, "mymark")
	if err != nil {
		t.Fatalf("ResolveBookmark: %v", err)
	}
	if id != "snap1" {
		t.Errorf("expected snap1, got %s", id)
	}
}

func TestAddBookmark_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := AddBookmark(dir, "x", "missing")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestResolveBookmark_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := ResolveBookmark(dir, "ghost")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRemoveBookmark(t *testing.T) {
	dir := t.TempDir()
	writeBookmarkSnapshot(t, dir, "snap2")
	_ = AddBookmark(dir, "mark2", "snap2")

	if err := RemoveBookmark(dir, "mark2"); err != nil {
		t.Fatalf("RemoveBookmark: %v", err)
	}
	_, err := ResolveBookmark(dir, "mark2")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestRemoveBookmark_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := RemoveBookmark(dir, "nope")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestListBookmarks(t *testing.T) {
	dir := t.TempDir()
	writeBookmarkSnapshot(t, dir, "s1")
	writeBookmarkSnapshot(t, dir, "s2")
	_ = AddBookmark(dir, "a", "s1")
	_ = AddBookmark(dir, "b", "s2")

	idx, err := ListBookmarks(dir)
	if err != nil {
		t.Fatalf("ListBookmarks: %v", err)
	}
	if len(idx) != 2 {
		t.Errorf("expected 2 bookmarks, got %d", len(idx))
	}
}
