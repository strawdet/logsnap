package snapshot_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryanfowler/logsnap/internal/snapshot"
)

func writeHighlightSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	snap := &snapshot.Snapshot{ID: id, Label: "hl-test"}
	data, _ := json.Marshal(snap)
	_ = os.WriteFile(filepath.Join(dir, id+".json"), data, 0644)
}

func TestAddHighlight_AndGet(t *testing.T) {
	dir := t.TempDir()
	writeHighlightSnapshot(t, dir, "snap1")

	if err := snapshot.AddHighlight(dir, "snap1", "error occurred"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	h, err := snapshot.GetHighlights(dir, "snap1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Messages) != 1 || h.Messages[0] != "error occurred" {
		t.Errorf("expected highlight, got %v", h.Messages)
	}
}

func TestAddHighlight_Deduplicates(t *testing.T) {
	dir := t.TempDir()
	writeHighlightSnapshot(t, dir, "snap2")

	_ = snapshot.AddHighlight(dir, "snap2", "dup msg")
	_ = snapshot.AddHighlight(dir, "snap2", "dup msg")

	h, _ := snapshot.GetHighlights(dir, "snap2")
	if len(h.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(h.Messages))
	}
}

func TestAddHighlight_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := snapshot.AddHighlight(dir, "missing", "msg")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestGetHighlights_NoFile(t *testing.T) {
	dir := t.TempDir()
	h, err := snapshot.GetHighlights(dir, "ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Messages) != 0 {
		t.Errorf("expected empty messages")
	}
}

func TestRemoveHighlight(t *testing.T) {
	dir := t.TempDir()
	writeHighlightSnapshot(t, dir, "snap3")
	_ = snapshot.AddHighlight(dir, "snap3", "keep")
	_ = snapshot.AddHighlight(dir, "snap3", "remove")

	if err := snapshot.RemoveHighlight(dir, "snap3", "remove"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	h, _ := snapshot.GetHighlights(dir, "snap3")
	for _, m := range h.Messages {
		if m == "remove" {
			t.Error("message should have been removed")
		}
	}
}

func TestRemoveHighlight_NotFound(t *testing.T) {
	dir := t.TempDir()
	writeHighlightSnapshot(t, dir, "snap4")
	_ = snapshot.AddHighlight(dir, "snap4", "only")

	err := snapshot.RemoveHighlight(dir, "snap4", "nonexistent")
	if err == nil {
		t.Error("expected error for missing highlight")
	}
}
