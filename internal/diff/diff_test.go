package diff_test

import (
	"strings"
	"testing"

	"github.com/user/logsnap/internal/diff"
	"github.com/user/logsnap/internal/snapshot"
)

func makeSnapshot(entries []snapshot.LogEntry) *snapshot.Snapshot {
	return &snapshot.Snapshot{Entries: entries}
}

var base = []snapshot.LogEntry{
	{Level: "INFO", Message: "server started"},
	{Level: "WARN", Message: "high memory usage"},
	{Level: "ERROR", Message: "db connection failed"},
}

func TestCompare_Added(t *testing.T) {
	current := []snapshot.LogEntry{
		{Level: "INFO", Message: "server started"},
		{Level: "WARN", Message: "high memory usage"},
		{Level: "ERROR", Message: "db connection failed"},
		{Level: "INFO", Message: "cache warmed up"},
	}

	result := diff.Compare(makeSnapshot(base), makeSnapshot(current))

	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added entry, got %d", len(result.Added))
	}
	if result.Added[0].Message != "cache warmed up" {
		t.Errorf("unexpected added message: %s", result.Added[0].Message)
	}
}

func TestCompare_Removed(t *testing.T) {
	current := []snapshot.LogEntry{
		{Level: "INFO", Message: "server started"},
	}

	result := diff.Compare(makeSnapshot(base), makeSnapshot(current))

	if len(result.Removed) != 2 {
		t.Fatalf("expected 2 removed entries, got %d", len(result.Removed))
	}
}

func TestCompare_Changed(t *testing.T) {
	current := []snapshot.LogEntry{
		{Level: "INFO", Message: "server started"},
		{Level: "ERROR", Message: "high memory usage"}, // level changed
		{Level: "ERROR", Message: "db connection failed"},
	}

	result := diff.Compare(makeSnapshot(base), makeSnapshot(current))

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed entry, got %d", len(result.Changed))
	}
	if result.Changed[0].OldLevel != "WARN" || result.Changed[0].NewLevel != "ERROR" {
		t.Errorf("unexpected level change: %+v", result.Changed[0])
	}
}

func TestResult_Summary(t *testing.T) {
	current := []snapshot.LogEntry{
		{Level: "INFO", Message: "server started"},
		{Level: "INFO", Message: "new feature enabled"},
	}

	result := diff.Compare(makeSnapshot(base), makeSnapshot(current))
	summary := result.Summary()

	if !strings.Contains(summary, "Added:") {
		t.Error("summary missing 'Added:' label")
	}
	if !strings.Contains(summary, "[+]") {
		t.Error("summary missing added entry marker")
	}
	if !strings.Contains(summary, "[-]") {
		t.Error("summary missing removed entry marker")
	}
}
