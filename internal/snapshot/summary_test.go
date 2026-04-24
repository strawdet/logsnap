package snapshot

import (
	"testing"
	"time"
)

func makeSummarySnapshot(dir string) *Snapshot {
	return &Snapshot{
		ID:        "sum-001",
		Label:     "summary-test",
		CreatedAt: time.Now(),
		Dir:       dir,
		Entries: []LogEntry{
			{Message: "disk full", Level: "error"},
			{Message: "disk full", Level: "error"},
			{Message: "disk full", Level: "error"},
			{Message: "connection timeout", Level: "warn"},
			{Message: "connection timeout", Level: "warn"},
			{Message: "service started", Level: "info"},
		},
	}
}

func TestSummarizeSnapshot_Counts(t *testing.T) {
	dir := t.TempDir()
	s := makeSummarySnapshot(dir)

	summary, err := SummarizeSnapshot(s, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Total != 6 {
		t.Errorf("expected total 6, got %d", summary.Total)
	}
	if summary.BySeverity["error"] != 3 {
		t.Errorf("expected 3 errors, got %d", summary.BySeverity["error"])
	}
	if summary.BySeverity["warn"] != 2 {
		t.Errorf("expected 2 warnings, got %d", summary.BySeverity["warn"])
	}
	if summary.BySeverity["info"] != 1 {
		t.Errorf("expected 1 info, got %d", summary.BySeverity["info"])
	}
}

func TestSummarizeSnapshot_TopMessages(t *testing.T) {
	dir := t.TempDir()
	s := makeSummarySnapshot(dir)

	summary, err := SummarizeSnapshot(s, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(summary.TopMessages) != 2 {
		t.Fatalf("expected 2 top messages, got %d", len(summary.TopMessages))
	}
	if summary.TopMessages[0] != "disk full" {
		t.Errorf("expected top message 'disk full', got %q", summary.TopMessages[0])
	}
}

func TestSummarizeSnapshot_NilSnapshot(t *testing.T) {
	_, err := SummarizeSnapshot(nil, 5)
	if err == nil {
		t.Error("expected error for nil snapshot, got nil")
	}
}

func TestSummarizeSnapshot_EmptyEntries(t *testing.T) {
	dir := t.TempDir()
	s := &Snapshot{
		ID:        "empty-001",
		Label:     "empty",
		CreatedAt: time.Now(),
		Dir:       dir,
		Entries:   []LogEntry{},
	}
	summary, err := SummarizeSnapshot(s, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Total != 0 {
		t.Errorf("expected total 0, got %d", summary.Total)
	}
	if len(summary.TopMessages) != 0 {
		t.Errorf("expected no top messages, got %d", len(summary.TopMessages))
	}
}
