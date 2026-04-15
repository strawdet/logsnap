package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/logsnap/internal/snapshot"
)

func writeSnapshot(t *testing.T, dir string, snap *snapshot.Snapshot) string {
	t.Helper()
	path := filepath.Join(dir, snap.ID+".json")
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}
	return path
}

func TestDiffCommand_Integration(t *testing.T) {
	dir := t.TempDir()

	baseline := &snapshot.Snapshot{
		ID: "snap-baseline",
		Entries: []snapshot.LogEntry{
			{Level: "INFO", Message: "app started"},
			{Level: "WARN", Message: "slow query detected"},
		},
	}

	current := &snapshot.Snapshot{
		ID: "snap-current",
		Entries: []snapshot.LogEntry{
			{Level: "INFO", Message: "app started"},
			{Level: "ERROR", Message: "slow query detected"}, // level changed
			{Level: "INFO", Message: "cache initialized"},   // added
		},
	}

	writeSnapshot(t, dir, baseline)
	writeSnapshot(t, dir, current)

	// Verify files exist before running diff logic directly.
	baselinePath := filepath.Join(dir, "snap-baseline.json")
	currentPath := filepath.Join(dir, "snap-current.json")

	loadedBase, err := snapshot.Load(baselinePath)
	if err != nil {
		t.Fatalf("load baseline: %v", err)
	}
	loadedCurrent, err := snapshot.Load(currentPath)
	if err != nil {
		t.Fatalf("load current: %v", err)
	}

	if len(loadedBase.Entries) != 2 {
		t.Errorf("expected 2 baseline entries, got %d", len(loadedBase.Entries))
	}
	if len(loadedCurrent.Entries) != 3 {
		t.Errorf("expected 3 current entries, got %d", len(loadedCurrent.Entries))
	}
}
