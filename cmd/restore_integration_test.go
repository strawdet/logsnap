package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nicholasgasior/logsnap/internal/snapshot"
)

func writeRestoreSnapshot(t *testing.T, dir string) string {
	t.Helper()
	snap, err := snapshot.New([]snapshot.LogEntry{
		{Timestamp: time.Now(), Level: "error", Message: "disk full"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := snap.Save(dir); err != nil {
		t.Fatal(err)
	}
	return snap.ID
}

func TestRestoreCommand_ToFile(t *testing.T) {
	dir := t.TempDir()
	id := writeRestoreSnapshot(t, dir)

	outFile := filepath.Join(dir, "out.log")

	result, err := snapshot.Restore(dir, id, outFile)
	if err != nil {
		t.Fatalf("restore error: %v", err)
	}
	if result.EntryCount != 1 {
		t.Errorf("expected 1 entry, got %d", result.EntryCount)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	var entry snapshot.LogEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if entry.Message != "disk full" {
		t.Errorf("unexpected message: %s", entry.Message)
	}
}
