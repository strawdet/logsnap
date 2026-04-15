package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTagSnapshot_AndResolve(t *testing.T) {
	dir := t.TempDir()

	err := TagSnapshot(dir, "production", "snap-abc123")
	if err != nil {
		t.Fatalf("TagSnapshot failed: %v", err)
	}

	id, err := ResolveTag(dir, "production")
	if err != nil {
		t.Fatalf("ResolveTag failed: %v", err)
	}
	if id != "snap-abc123" {
		t.Errorf("expected snap-abc123, got %s", id)
	}
}

func TestResolveTag_NotFound(t *testing.T) {
	dir := t.TempDir()

	_, err := ResolveTag(dir, "missing")
	if err == nil {
		t.Fatal("expected error for missing tag, got nil")
	}
}

func TestRemoveTag(t *testing.T) {
	dir := t.TempDir()

	_ = TagSnapshot(dir, "staging", "snap-xyz789")

	err := RemoveTag(dir, "staging")
	if err != nil {
		t.Fatalf("RemoveTag failed: %v", err)
	}

	_, err = ResolveTag(dir, "staging")
	if err == nil {
		t.Fatal("expected error after removal, got nil")
	}
}

func TestRemoveTag_NotFound(t *testing.T) {
	dir := t.TempDir()

	err := RemoveTag(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent tag")
	}
}

func TestLoadTagIndex_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	index, err := LoadTagIndex(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(index) != 0 {
		t.Errorf("expected empty index, got %d entries", len(index))
	}
}

func TestLoadTagIndex_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "tags.json"), []byte("not-json{"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = LoadTagIndex(dir)
	if err == nil {
		t.Fatal("expected JSON parse error, got nil")
	}
}

func TestTagSnapshot_Overwrite(t *testing.T) {
	dir := t.TempDir()

	_ = TagSnapshot(dir, "v1", "snap-old")
	_ = TagSnapshot(dir, "v1", "snap-new")

	id, err := ResolveTag(dir, "v1")
	if err != nil {
		t.Fatalf("ResolveTag failed: %v", err)
	}
	if id != "snap-new" {
		t.Errorf("expected snap-new after overwrite, got %s", id)
	}
}
