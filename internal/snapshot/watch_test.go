package snapshot

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeLogFile(t *testing.T, path string, entries []LogEntry) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create log file: %v", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			t.Fatalf("encode entry: %v", err)
		}
	}
}

func TestWatch_CapturesSnapshotOnChange(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "app.log")

	entries := []LogEntry{
		{Level: "info", Message: "started", Timestamp: "2024-01-01T00:00:00Z"},
		{Level: "error", Message: "oops", Timestamp: "2024-01-01T00:00:01Z"},
	}
	writeLogFile(t, logFile, entries)

	snapped := make(chan string, 1)
	opts := WatchOptions{
		LogFile:  logFile,
		Dir:      dir,
		Interval: 50 * time.Millisecond,
		OnSnap: func(id string, err error) {
			if err == nil {
				snapped <- id
			}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go Watch(ctx, opts) //nolint:errcheck

	select {
	case id := <-snapped:
		if id == "" {
			t.Fatal("expected non-empty snapshot ID")
		}
		snap, err := Load(dir, id)
		if err != nil {
			t.Fatalf("load snapshot: %v", err)
		}
		if len(snap.Entries) != 2 {
			t.Errorf("expected 2 entries, got %d", len(snap.Entries))
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for snapshot")
	}
}

func TestWatch_SkipsUnchangedFile(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "app.log")
	writeLogFile(t, logFile, []LogEntry{
		{Level: "info", Message: "hello", Timestamp: "2024-01-01T00:00:00Z"},
	})

	count := 0
	opts := WatchOptions{
		LogFile:  logFile,
		Dir:      dir,
		Interval: 40 * time.Millisecond,
		OnSnap: func(_ string, err error) {
			if err == nil {
				count++
			}
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	Watch(ctx, opts) //nolint:errcheck

	if count > 1 {
		t.Errorf("expected at most 1 snapshot for unchanged file, got %d", count)
	}
}
