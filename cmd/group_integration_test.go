package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"logsnap/internal/snapshot"
)

func writeGroupSnap(t *testing.T, dir, id string) {
	t.Helper()
	s := &snapshot.Snapshot{ID: id, Label: "group-test"}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, id+".json"), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestGroupCommand_AddAndShow(t *testing.T) {
	dir := t.TempDir()
	writeGroupSnap(t, dir, "snap-aaa")
	writeGroupSnap(t, dir, "snap-bbb")

	root := rootCmdWithGroup()

	root.SetArgs([]string{"group", "add", "my-group", "snap-aaa", "--dir", dir})
	if err := root.Execute(); err != nil {
		t.Fatalf("group add: %v", err)
	}

	root.SetArgs([]string{"group", "add", "my-group", "snap-bbb", "--dir", dir})
	if err := root.Execute(); err != nil {
		t.Fatalf("group add second: %v", err)
	}

	g, err := snapshot.GetGroup(dir, "my-group")
	if err != nil {
		t.Fatalf("GetGroup: %v", err)
	}
	if len(g.Snapshots) != 2 {
		t.Errorf("expected 2 snapshots in group, got %d", len(g.Snapshots))
	}
}

func TestGroupCommand_List(t *testing.T) {
	dir := t.TempDir()
	writeGroupSnap(t, dir, "snap-ccc")

	root := rootCmdWithGroup()
	root.SetArgs([]string{"group", "add", "g1", "snap-ccc", "--dir", dir})
	_ = root.Execute()

	root.SetArgs([]string{"group", "add", "g2", "snap-ccc", "--dir", dir})
	_ = root.Execute()

	names, err := snapshot.ListGroups(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 groups, got %d", len(names))
	}
}

func rootCmdWithGroup() *cobra.Command {
	return RootCmd
}
