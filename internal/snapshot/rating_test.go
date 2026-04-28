package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeRatingSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	snap := &Snapshot{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Entries:   []LogEntry{},
	}
	data, _ := json.Marshal(snap)
	_ = os.WriteFile(filepath.Join(dir, id+".json"), data, 0644)
}

func TestSetRating_AndGet(t *testing.T) {
	dir := t.TempDir()
	writeRatingSnapshot(t, dir, "snap1")

	if err := SetRating(dir, "snap1", 4, "looks good"); err != nil {
		t.Fatalf("SetRating: %v", err)
	}

	r, err := GetRating(dir, "snap1")
	if err != nil {
		t.Fatalf("GetRating: %v", err)
	}
	if r.Score != 4 {
		t.Errorf("expected score 4, got %d", r.Score)
	}
	if r.Comment != "looks good" {
		t.Errorf("expected comment 'looks good', got %q", r.Comment)
	}
}

func TestSetRating_PreservesCreatedAt(t *testing.T) {
	dir := t.TempDir()
	writeRatingSnapshot(t, dir, "snap2")

	_ = SetRating(dir, "snap2", 3, "initial")
	first, _ := GetRating(dir, "snap2")

	_ = SetRating(dir, "snap2", 5, "updated")
	second, _ := GetRating(dir, "snap2")

	if !second.CreatedAt.Equal(first.CreatedAt) {
		t.Errorf("CreatedAt should be preserved on update")
	}
	if second.Score != 5 {
		t.Errorf("expected updated score 5, got %d", second.Score)
	}
}

func TestSetRating_InvalidScore(t *testing.T) {
	dir := t.TempDir()
	writeRatingSnapshot(t, dir, "snap3")

	for _, score := range []int{0, 6, -1} {
		if err := SetRating(dir, "snap3", score, ""); err == nil {
			t.Errorf("expected error for score %d", score)
		}
	}
}

func TestSetRating_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	if err := SetRating(dir, "missing", 3, "hi"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestGetRating_NoFile(t *testing.T) {
	dir := t.TempDir()
	writeRatingSnapshot(t, dir, "snap4")

	_, err := GetRating(dir, "snap4")
	if err == nil {
		t.Error("expected error when no rating file exists")
	}
}

func TestRemoveRating(t *testing.T) {
	dir := t.TempDir()
	writeRatingSnapshot(t, dir, "snap5")
	_ = SetRating(dir, "snap5", 2, "meh")

	if err := RemoveRating(dir, "snap5"); err != nil {
		t.Fatalf("RemoveRating: %v", err)
	}
	if _, err := GetRating(dir, "snap5"); err == nil {
		t.Error("expected error after removal")
	}
}

func TestRemoveRating_NotFound(t *testing.T) {
	dir := t.TempDir()
	if err := RemoveRating(dir, "ghost"); err == nil {
		t.Error("expected error removing non-existent rating")
	}
}
