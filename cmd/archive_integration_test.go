package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"logsnap/internal/snapshot"
)

func writeArchiveSnap(t *testing.T, dir, id string) {
	t.Helper()
	snap := &snapshot.Snapshot{
		ID:        id,
		Label:     "archive-test",
		CreatedAt: time.Now(),
		Entries:   []snapshot.LogEntry{{Level: "warn", Message: "test", Timestamp: time.Now()}},
	}
	data, _ := json.Marshal(snap)
	_ = os.WriteFile(filepath.Join(dir, id+".json"), data, 0644)
}

func TestArchiveCommand_RoundTrip(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()
	archivePath := filepath.Join(t.TempDir(), "snaps.zip")

	writeArchiveSnap(t, srcDir, "arc-001")
	writeArchiveSnap(t, srcDir, "arc-002")

	root := rootCmd
	root.SetArgs([]string{"archive", "--dir", srcDir, "--out", archivePath, "arc-001", "arc-002"})
	if err := root.Execute(); err != nil {
		t.Fatalf("archive command: %v", err)
	}

	if _, err := os.Stat(archivePath); err != nil {
		t.Fatalf("archive not created: %v", err)
	}

	root.SetArgs([]string{"unarchive", "--dir", dstDir, archivePath})
	if err := root.Execute(); err != nil {
		t.Fatalf("unarchive command: %v", err)
	}

	for _, id := range []string{"arc-001", "arc-002"} {
		if _, err := os.Stat(filepath.Join(dstDir, id+".json")); err != nil {
			t.Errorf("restored snapshot %s missing", id)
		}
	}
}
