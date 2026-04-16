package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"logsnap/internal/snapshot"
)

func writeAliasSnap(t *testing.T, dir, id string) {
	t.Helper()
	f, err := os.Create(filepath.Join(dir, id+".json"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func TestAliasCommand_SetAndGet(t *testing.T) {
	dir := t.TempDir()
	writeAliasSnap(t, dir, "snap-001")

	if err := snapshot.SetAlias(dir, "live", "snap-001"); err != nil {
		t.Fatalf("SetAlias: %v", err)
	}

	id, err := snapshot.ResolveAlias(dir, "live")
	if err != nil {
		t.Fatalf("ResolveAlias: %v", err)
	}
	if id != "snap-001" {
		t.Errorf("expected snap-001, got %s", id)
	}
}

func TestAliasCommand_List(t *testing.T) {
	dir := t.TempDir()
	writeAliasSnap(t, dir, "snap-A")
	writeAliasSnap(t, dir, "snap-B")
	_ = snapshot.SetAlias(dir, "one", "snap-A")
	_ = snapshot.SetAlias(dir, "two", "snap-B")

	idx, err := snapshot.ListAliases(dir)
	if err != nil {
		t.Fatalf("ListAliases: %v", err)
	}
	if len(idx) != 2 {
		t.Errorf("expected 2, got %d", len(idx))
	}
	for _, key := range []string{"one", "two"} {
		if _, ok := idx[key]; !ok {
			t.Errorf("missing alias %s", key)
		}
	}
}

func TestAliasCommand_Remove(t *testing.T) {
	dir := t.TempDir()
	writeAliasSnap(t, dir, "snap-Z")
	_ = snapshot.SetAlias(dir, "temp", "snap-Z")

	if err := snapshot.RemoveAlias(dir, "temp"); err != nil {
		t.Fatalf("RemoveAlias: %v", err)
	}
	idx, _ := snapshot.ListAliases(dir)
	if _, ok := idx["temp"]; ok {
		t.Error("alias should have been removed")
	}
}

func rootCmdWithAlias() *cobra.Command { _ = strings.Join(nil, ""); return RootCmd }
