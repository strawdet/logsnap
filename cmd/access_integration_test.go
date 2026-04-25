package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"logsnap/internal/snapshot"
)

func writeAccessSnap(t *testing.T, dir, id string) {
	t.Helper()
	f := filepath.Join(dir, id+".json")
	if err := os.WriteFile(f, []byte(`{"id":"`+id+`","entries":[]}`), 0644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}
}

func rootCmdWithAccess() *cobra.Command {
	root := &cobra.Command{Use: "logsnap"}
	accessCmd := &cobra.Command{Use: "access"}

	var dir string
	accessCmd.PersistentFlags().StringVar(&dir, "dir", ".logsnap", "")

	recordCmd := &cobra.Command{
		Use:  "record",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			action, _ := cmd.Flags().GetString("action")
			actor, _ := cmd.Flags().GetString("actor")
			note, _ := cmd.Flags().GetString("note")
			return snapshot.RecordAccess(dir, args[0], action, actor, note)
		},
	}
	recordCmd.Flags().String("action", "", "")
	recordCmd.Flags().String("actor", "", "")
	recordCmd.Flags().String("note", "", "")

	showCmd := &cobra.Command{
		Use:  "show",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			log, err := snapshot.GetAccessLog(dir, args[0])
			if err != nil {
				return err
			}
			for _, e := range log.Events {
				cmd.Printf("%s %s %s\n", e.Action, e.Actor, e.Note)
			}
			return nil
		},
	}

	accessCmd.AddCommand(recordCmd, showCmd)
	root.AddCommand(accessCmd)
	return root
}

func TestAccessCommand_RecordAndShow(t *testing.T) {
	dir := t.TempDir()
	writeAccessSnap(t, dir, "snap42")

	root := rootCmdWithAccess()

	root.SetArgs([]string{"access", "record", "snap42",
		"--dir", dir, "--action", "read", "--actor", "tester", "--note", "ci"})
	if err := root.Execute(); err != nil {
		t.Fatalf("record: %v", err)
	}

	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"access", "show", "snap42", "--dir", dir})
	if err := root.Execute(); err != nil {
		t.Fatalf("show: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "read") || !strings.Contains(out, "tester") {
		t.Errorf("unexpected output: %q", out)
	}
}
