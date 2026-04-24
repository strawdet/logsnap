package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func writeNamespaceSnap(t *testing.T, dir, id string) {
	t.Helper()
	snap := &snapshot.Snapshot{
		ID:      id,
		Entries: []snapshot.LogEntry{{Level: "info", Message: "ns integration"}},
	}
	data, _ := json.Marshal(snap)
	_ = os.WriteFile(filepath.Join(dir, id+".json"), data, 0644)
}

func rootCmdWithNamespace(dir string) *cobra.Command {
	root := &cobra.Command{Use: "logsnap"}
	ns := &cobra.Command{Use: "namespace"}

	add := &cobra.Command{
		Use:  "add",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return snapshot.AddToNamespace(dir, args[0], args[1])
		},
	}
	list := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			names, err := snapshot.ListNamespaces(dir)
			if err != nil {
				return err
			}
			for _, n := range names {
				cmd.Println(n)
			}
			return nil
		},
	}
	show := &cobra.Command{
		Use:  "show",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, err := snapshot.GetNamespaceSnapshots(dir, args[0])
			if err != nil {
				return err
			}
			for _, id := range ids {
				cmd.Println(id)
			}
			return nil
		},
	}
	ns.AddCommand(add, list, show)
	root.AddCommand(ns)
	return root
}

func TestNamespaceCommand_AddAndList(t *testing.T) {
	dir := t.TempDir()
	writeNamespaceSnap(t, dir, "snap-int-1")
	writeNamespaceSnap(t, dir, "snap-int-2")

	root := rootCmdWithNamespace(dir)

	root.SetArgs([]string{"namespace", "add", "production", "snap-int-1"})
	if err := root.Execute(); err != nil {
		t.Fatalf("add: %v", err)
	}
	root.SetArgs([]string{"namespace", "add", "production", "snap-int-2"})
	if err := root.Execute(); err != nil {
		t.Fatalf("add: %v", err)
	}

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"namespace", "list"})
	if err := root.Execute(); err != nil {
		t.Fatalf("list: %v", err)
	}
	if !strings.Contains(buf.String(), "production") {
		t.Errorf("expected 'production' in output, got: %s", buf.String())
	}
}

func TestNamespaceCommand_Show(t *testing.T) {
	dir := t.TempDir()
	writeNamespaceSnap(t, dir, "snap-show-1")

	root := rootCmdWithNamespace(dir)
	root.SetArgs([]string{"namespace", "add", "dev", "snap-show-1"})
	_ = root.Execute()

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"namespace", "show", "dev"})
	if err := root.Execute(); err != nil {
		t.Fatalf("show: %v", err)
	}
	if !strings.Contains(buf.String(), "snap-show-1") {
		t.Errorf("expected snap-show-1 in output, got: %s", buf.String())
	}
}
