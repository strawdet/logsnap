package snapshot_test

import (
	"os"
	"testing"

	"github.com/yourusername/logsnap/internal/snapshot"
)

func writeBadgeSnapshot(t *testing.T, dir string) *snapshot.Snapshot {
	t.Helper()
	entries := []snapshot.LogEntry{
		{Level: "info", Message: "badge test"},
	}
	snap, err := snapshot.New(dir, entries, "badge-test")
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}
	return snap
}

func TestAddBadge_AndGet(t *testing.T) {
	dir := t.TempDir()
	snap := writeBadgeSnapshot(t, dir)

	err := snapshot.AddBadge(dir, snap.ID, "gold", "🥇", "top performer")
	if err != nil {
		t.Fatalf("AddBadge failed: %v", err)
	}

	index, err := snapshot.GetBadges(dir, snap.ID)
	if err != nil {
		t.Fatalf("GetBadges failed: %v", err)
	}
	if len(index.Badges) != 1 {
		t.Fatalf("expected 1 badge, got %d", len(index.Badges))
	}
	if index.Badges[0].Name != "gold" {
		t.Errorf("expected badge name 'gold', got %q", index.Badges[0].Name)
	}
	if index.Badges[0].Icon != "🥇" {
		t.Errorf("expected icon '🥇', got %q", index.Badges[0].Icon)
	}
}

func TestAddBadge_Deduplicates(t *testing.T) {
	dir := t.TempDir()
	snap := writeBadgeSnapshot(t, dir)

	_ = snapshot.AddBadge(dir, snap.ID, "silver", "🥈", "first try")
	_ = snapshot.AddBadge(dir, snap.ID, "silver", "🥈", "duplicate")

	index, _ := snapshot.GetBadges(dir, snap.ID)
	if len(index.Badges) != 1 {
		t.Errorf("expected 1 badge after dedup, got %d", len(index.Badges))
	}
}

func TestAddBadge_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := snapshot.AddBadge(dir, "nonexistent-id", "gold", "🥇", "reason")
	if err == nil {
		t.Error("expected error for missing snapshot, got nil")
	}
}

func TestGetBadges_NoFile(t *testing.T) {
	dir := t.TempDir()
	index, err := snapshot.GetBadges(dir, "any-id")
	if err != nil {
		t.Fatalf("expected no error for missing badge file, got %v", err)
	}
	if len(index.Badges) != 0 {
		t.Errorf("expected empty badges, got %d", len(index.Badges))
	}
}

func TestRemoveBadge(t *testing.T) {
	dir := t.TempDir()
	snap := writeBadgeSnapshot(t, dir)

	_ = snapshot.AddBadge(dir, snap.ID, "bronze", "🥉", "participation")
	err := snapshot.RemoveBadge(dir, snap.ID, "bronze")
	if err != nil {
		t.Fatalf("RemoveBadge failed: %v", err)
	}

	index, _ := snapshot.GetBadges(dir, snap.ID)
	if len(index.Badges) != 0 {
		t.Errorf("expected 0 badges after removal, got %d", len(index.Badges))
	}
}

func TestRemoveBadge_NotFound(t *testing.T) {
	dir := t.TempDir()
	snap := writeBadgeSnapshot(t, dir)

	err := snapshot.RemoveBadge(dir, snap.ID, "nonexistent")
	if err == nil {
		t.Error("expected error when removing nonexistent badge")
	}
	_ = os.RemoveAll(dir)
}
