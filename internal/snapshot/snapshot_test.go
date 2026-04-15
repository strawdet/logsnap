package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/logsnap/logsnap/internal/snapshot"
)

func sampleEntries() []snapshot.LogEntry {
	return []snapshot.LogEntry{
		{
			Timestamp: "2024-01-01T00:00:00Z",
			Level:     "INFO",
			Message:   "service started",
			Fields:    map[string]string{"service": "api"},
		},
		{
			Timestamp: "2024-01-01T00:00:01Z",
			Level:     "ERROR",
			Message:   "connection refused",
			Fields:    map[string]string{"host": "db:5432"},
		},
	}
}

func TestNew(t *testing.T) {
	s := snapshot.New("v1.0.0", sampleEntries(), nil)
	if s.Label != "v1.0.0" {
		t.Errorf("expected label v1.0.0, got %s", s.Label)
	}
	if s.ID == "" {
		t.Error("expected non-empty ID")
	}
	if len(s.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(s.Entries))
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	orig := snapshot.New("v1.0.0", sampleEntries(), map[string]string{"env": "prod"})
	if err := orig.Save(path); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.ID != orig.ID {
		t.Errorf("ID mismatch: got %s, want %s", loaded.ID, orig.ID)
	}
	if loaded.Label != orig.Label {
		t.Errorf("Label mismatch: got %s, want %s", loaded.Label, orig.Label)
	}
	if len(loaded.Entries) != len(orig.Entries) {
		t.Errorf("Entries count mismatch: got %d, want %d", len(loaded.Entries), len(orig.Entries))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
