package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func writeAliasSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	f, err := os.Create(filepath.Join(dir, id+".json"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func TestSetAlias_AndResolve(t *testing.T) {
	dir := t.TempDir()
	writeAliasSnapshot(t, dir, "snap-abc")
	if err := SetAlias(dir, "prod", "snap-abc"); err != nil {
		t.Fatalf("SetAlias: %v", err)
	}
	id, err := ResolveAlias(dir, "prod")
	if err != nil {
		t.Fatalf("ResolveAlias: %v", err)
	}
	if id != "snap-abc" {
		t.Errorf("expected snap-abc, got %s", id)
	}
}

func TestSetAlias_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	if err := SetAlias(dir, "prod", "missing"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestResolveAlias_NotFound(t *testing.T) {
	dir := t.TempDir()
	if _, err := ResolveAlias(dir, "ghost"); err == nil {
		t.Error("expected error for missing alias")
	}
}

func TestRemoveAlias(t *testing.T) {
	dir := t.TempDir()
	writeAliasSnapshot(t, dir, "snap-xyz")
	_ = SetAlias(dir, "staging", "snap-xyz")
	if err := RemoveAlias(dir, "staging"); err != nil {
		t.Fatalf("RemoveAlias: %v", err)
	}
	if _, err := ResolveAlias(dir, "staging"); err == nil {
		t.Error("expected alias to be removed")
	}
}

func TestRemoveAlias_NotFound(t *testing.T) {
	dir := t.TempDir()
	if err := RemoveAlias(dir, "nope"); err == nil {
		t.Error("expected error for missing alias")
	}
}

func TestListAliases(t *testing.T) {
	dir := t.TempDir()
	writeAliasSnapshot(t, dir, "snap-1")
	writeAliasSnapshot(t, dir, "snap-2")
	_ = SetAlias(dir, "alpha", "snap-1")
	_ = SetAlias(dir, "beta", "snap-2")
	idx, err := ListAliases(dir)
	if err != nil {
		t.Fatalf("ListAliases: %v", err)
	}
	if len(idx) != 2 {
		t.Errorf("expected 2 aliases, got %d", len(idx))
	}
}
