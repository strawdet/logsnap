package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTimelineSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	snap := &Snapshot{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Entries:   []LogEntry{},
	}
	data, _ := json.Marshal(snap)
	_ = os.WriteFile(filepath.Join(dir, id+".json"), data, 0644)
}

func TestAddTimelineEvent_AndGet(t *testing.T) {
	dir := t.TempDir()
	writeTimelineSnapshot(t, dir, "snap1")

	if err := AddTimelineEvent(dir, "snap1", "captured", "initial capture"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := AddTimelineEvent(dir, "snap1", "tagged", "tag: prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tl, err := GetTimeline(dir, "snap1")
	if err != nil {
		t.Fatalf("GetTimeline error: %v", err)
	}
	if len(tl.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(tl.Entries))
	}
	if tl.Entries[0].Event != "captured" {
		t.Errorf("expected first event 'captured', got %q", tl.Entries[0].Event)
	}
	if tl.Entries[1].Detail != "tag: prod" {
		t.Errorf("expected detail 'tag: prod', got %q", tl.Entries[1].Detail)
	}
}

func TestGetTimeline_NoFile(t *testing.T) {
	dir := t.TempDir()
	writeTimelineSnapshot(t, dir, "snap2")

	tl, err := GetTimeline(dir, "snap2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tl.Entries) != 0 {
		t.Errorf("expected empty timeline, got %d entries", len(tl.Entries))
	}
}

func TestAddTimelineEvent_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := AddTimelineEvent(dir, "ghost", "captured", "")
	if err == nil {
		t.Fatal("expected error for missing snapshot, got nil")
	}
}

func TestClearTimeline(t *testing.T) {
	dir := t.TempDir()
	writeTimelineSnapshot(t, dir, "snap3")

	_ = AddTimelineEvent(dir, "snap3", "captured", "")

	if err := ClearTimeline(dir, "snap3"); err != nil {
		t.Fatalf("ClearTimeline error: %v", err)
	}

	tl, err := GetTimeline(dir, "snap3")
	if err != nil {
		t.Fatalf("GetTimeline after clear: %v", err)
	}
	if len(tl.Entries) != 0 {
		t.Errorf("expected empty timeline after clear, got %d entries", len(tl.Entries))
	}
}

func TestClearTimeline_NoFile(t *testing.T) {
	dir := t.TempDir()
	if err := ClearTimeline(dir, "nosnap"); err != nil {
		t.Errorf("expected no error clearing non-existent timeline, got: %v", err)
	}
}
