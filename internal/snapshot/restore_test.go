package snapshot_test

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nicholasgasior/logsnap/internal/snapshot"
)

func makeRestoreSnapshot(t *testing.T, dir string) *snapshot.Snapshot {
	t.Helper()
	snap, err := snapshot.New([]snapshot.LogEntry{
		{Timestamp: time.Now(), Level: "info", Message: "restore me"},
		{Timestamp: time.Now(), Level: "warn", Message: "and me"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := snap.Save(dir); err != nil {
		t.Fatal(err)
	}
	return snap
}

func TestRestore_ToFile(t *testing.T) {
	dir := t.TempDir()
	snap := makeRestoreSnapshot(t, dir)

	out := filepath.Join(dir, "restored.log")
	res, err := snapshot.Restore(dir, snap.ID, out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.EntryCount != 2 {
		t.Errorf("expected 2 entries, got %d", res.EntryCount)
	}

	f, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var entries []snapshot.LogEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var e snapshot.LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
			t.Fatal(err)
		}
		entries = append(entries, e)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 decoded entries, got %d", len(entries))
	}
}

func TestRestore_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := snapshot.Restore(dir, "nonexistent", "-")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}
