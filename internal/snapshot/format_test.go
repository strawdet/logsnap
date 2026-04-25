package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/logsnap/internal/snapshot"
)

func writeFormatSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	snap := &snapshot.Snapshot{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Entries:   []snapshot.LogEntry{{Level: "info", Message: "boot"}},
	}
	if err := snap.Save(dir); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
}

func TestSetFormatConfig_AndGet(t *testing.T) {
	dir := t.TempDir()
	writeFormatSnapshot(t, dir, "snap1")

	cfg := snapshot.FormatConfig{
		TimestampLayout: time.RFC3339,
		DateFormat:      "2006-01-02",
		LevelColors:     map[string]string{"error": "red", "info": "green"},
	}
	if err := snapshot.SetFormatConfig(dir, "snap1", cfg); err != nil {
		t.Fatalf("SetFormatConfig: %v", err)
	}

	got, err := snapshot.GetFormatConfig(dir, "snap1")
	if err != nil {
		t.Fatalf("GetFormatConfig: %v", err)
	}
	if got.TimestampLayout != time.RFC3339 {
		t.Errorf("TimestampLayout: got %q, want %q", got.TimestampLayout, time.RFC3339)
	}
	if got.LevelColors["error"] != "red" {
		t.Errorf("LevelColors[error]: got %q, want %q", got.LevelColors["error"], "red")
	}
	if got.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if got.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

func TestSetFormatConfig_PreservesCreatedAt(t *testing.T) {
	dir := t.TempDir()
	writeFormatSnapshot(t, dir, "snap2")

	cfg := snapshot.FormatConfig{TimestampLayout: "first"}
	if err := snapshot.SetFormatConfig(dir, "snap2", cfg); err != nil {
		t.Fatalf("first SetFormatConfig: %v", err)
	}
	first, _ := snapshot.GetFormatConfig(dir, "snap2")

	cfg2 := snapshot.FormatConfig{TimestampLayout: "second"}
	if err := snapshot.SetFormatConfig(dir, "snap2", cfg2); err != nil {
		t.Fatalf("second SetFormatConfig: %v", err)
	}
	second, _ := snapshot.GetFormatConfig(dir, "snap2")

	if !second.CreatedAt.Equal(first.CreatedAt) {
		t.Errorf("CreatedAt changed on update: got %v, want %v", second.CreatedAt, first.CreatedAt)
	}
}

func TestSetFormatConfig_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	cfg := snapshot.FormatConfig{TimestampLayout: time.RFC3339}
	err := snapshot.SetFormatConfig(dir, "ghost", cfg)
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestGetFormatConfig_NoFile(t *testing.T) {
	dir := t.TempDir()
	writeFormatSnapshot(t, dir, "snap3")
	_, err := snapshot.GetFormatConfig(dir, "snap3")
	if err == nil {
		t.Fatal("expected error when no format config file exists")
	}
}

func TestRemoveFormatConfig(t *testing.T) {
	dir := t.TempDir()
	writeFormatSnapshot(t, dir, "snap4")
	cfg := snapshot.FormatConfig{DateFormat: "2006"}
	if err := snapshot.SetFormatConfig(dir, "snap4", cfg); err != nil {
		t.Fatalf("SetFormatConfig: %v", err)
	}
	if err := snapshot.RemoveFormatConfig(dir, "snap4"); err != nil {
		t.Fatalf("RemoveFormatConfig: %v", err)
	}
	path := filepath.Join(dir, "snap4.format.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected format config file to be deleted")
	}
}

func TestRemoveFormatConfig_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := snapshot.RemoveFormatConfig(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error when removing non-existent format config")
	}
}
