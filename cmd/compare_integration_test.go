package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logsnap/internal/snapshot"
)

func writeCompareSnapshot(t *testing.T, dir, id string, entries []snapshot.LogEntry) {
	t.Helper()
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	snap := &snapshot.Snapshot{
		ID:        id,
		CreatedAt: time.Now(),
		Entries:   entries,
	}
	data, _ := json.MarshalIndent(snap, "", "  ")
	path := filepath.Join(dir, id+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}
}

func TestCompareCommand_SaveAndList(t *testing.T) {
	snapshotDir := t.TempDir()
	compareDir := t.TempDir()

	baseEntries := []snapshot.LogEntry{
		{Level: "info", Message: "server started", Timestamp: time.Now()},
		{Level: "warn", Message: "high memory", Timestamp: time.Now()},
	}
	targetEntries := []snapshot.LogEntry{
		{Level: "info", Message: "server started", Timestamp: time.Now()},
		{Level: "error", Message: "disk full", Timestamp: time.Now()},
	}

	writeCompareSnapshot(t, snapshotDir, "base001", baseEntries)
	writeCompareSnapshot(t, snapshotDir, "target001", targetEntries)

	_, err := snapshot.SaveCompareResult(compareDir, "release-1.0", "base001", "target001")
	if err != nil {
		t.Fatalf("SaveCompareResult: %v", err)
	}

	names, err := snapshot.ListCompareResults(compareDir)
	if err != nil {
		t.Fatalf("ListCompareResults: %v", err)
	}

	if len(names) != 1 || names[0] != "release-1.0" {
		t.Errorf("expected [release-1.0], got %v", names)
	}

	loaded, err := snapshot.LoadCompareResult(compareDir, "release-1.0")
	if err != nil {
		t.Fatalf("LoadCompareResult: %v", err)
	}
	if loaded.BaseID != "base001" || loaded.TargetID != "target001" {
		t.Errorf("unexpected IDs: base=%s target=%s", loaded.BaseID, loaded.TargetID)
	}
}
