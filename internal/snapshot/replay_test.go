package snapshot

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func makeReplaySnapshot() *Snapshot {
	now := time.Now()
	return &Snapshot{
		ID:        "replay-test-id",
		CreatedAt: now,
		Entries: []LogEntry{
			{Timestamp: now, Level: "info", Message: "service started"},
			{Timestamp: now, Level: "error", Message: "connection failed"},
			{Timestamp: now, Level: "info", Message: "retrying connection"},
			{Timestamp: now, Level: "error", Message: "timeout exceeded"},
		},
	}
}

func TestReplay_AllEntries(t *testing.T) {
	snap := makeReplaySnapshot()
	var buf bytes.Buffer

	result, err := Replay(snap, ReplayOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Replayed != 4 {
		t.Errorf("expected 4 replayed, got %d", result.Replayed)
	}
	if result.Skipped != 0 {
		t.Errorf("expected 0 skipped, got %d", result.Skipped)
	}
	if !strings.Contains(buf.String(), "service started") {
		t.Error("expected output to contain 'service started'")
	}
}

func TestReplay_WithFilter(t *testing.T) {
	snap := makeReplaySnapshot()
	var buf bytes.Buffer

	result, err := Replay(snap, ReplayOptions{Writer: &buf, Filter: "error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Replayed != 2 {
		t.Errorf("expected 2 replayed, got %d", result.Replayed)
	}
	if result.Skipped != 2 {
		t.Errorf("expected 2 skipped, got %d", result.Skipped)
	}
	if strings.Contains(buf.String(), "service started") {
		t.Error("expected 'service started' to be filtered out")
	}
}

func TestReplay_NilSnapshot(t *testing.T) {
	_, err := Replay(nil, ReplayOptions{})
	if err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestReplay_EmptyEntries(t *testing.T) {
	snap := &Snapshot{ID: "empty", Entries: []LogEntry{}}
	var buf bytes.Buffer

	result, err := Replay(snap, ReplayOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Replayed != 0 {
		t.Errorf("expected 0 replayed, got %d", result.Replayed)
	}
}
