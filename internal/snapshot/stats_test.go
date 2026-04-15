package snapshot

import (
	"testing"
	"time"
)

func makeStatsSnapshot() *Snapshot {
	return &Snapshot{
		ID:        "test-stats-id",
		Label:     "stats-test",
		CreatedAt: time.Now(),
		Entries: []LogEntry{
			{Level: "error", Message: "disk full"},
			{Level: "error", Message: "disk full"},
			{Level: "warn", Message: "high memory"},
			{Level: "info", Message: "server started"},
			{Level: "info", Message: "server started"},
			{Level: "info", Message: "server started"},
			{Level: "", Message: "unknown origin"},
		},
	}
}

func TestComputeStats_Counts(t *testing.T) {
	snap := makeStatsSnapshot()
	stats, err := ComputeStats(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.TotalCount != 7 {
		t.Errorf("expected TotalCount 7, got %d", stats.TotalCount)
	}
	if stats.LevelCounts["error"] != 2 {
		t.Errorf("expected 2 errors, got %d", stats.LevelCounts["error"])
	}
	if stats.LevelCounts["info"] != 3 {
		t.Errorf("expected 3 infos, got %d", stats.LevelCounts["info"])
	}
	if stats.LevelCounts["unknown"] != 1 {
		t.Errorf("expected 1 unknown, got %d", stats.LevelCounts["unknown"])
	}
}

func TestComputeStats_TopMessages(t *testing.T) {
	snap := makeStatsSnapshot()
	stats, err := ComputeStats(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stats.TopMessages) == 0 {
		t.Fatal("expected non-empty TopMessages")
	}
	if stats.TopMessages[0] != "server started" {
		t.Errorf("expected top message 'server started', got %q", stats.TopMessages[0])
	}
}

func TestComputeStats_NilSnapshot(t *testing.T) {
	_, err := ComputeStats(nil)
	if err == nil {
		t.Error("expected error for nil snapshot, got nil")
	}
}

func TestComputeStats_EmptyEntries(t *testing.T) {
	snap := &Snapshot{ID: "empty", Label: "none", Entries: []LogEntry{}}
	stats, err := ComputeStats(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.TotalCount != 0 {
		t.Errorf("expected TotalCount 0, got %d", stats.TotalCount)
	}
	if len(stats.TopMessages) != 0 {
		t.Errorf("expected empty TopMessages, got %v", stats.TopMessages)
	}
}
