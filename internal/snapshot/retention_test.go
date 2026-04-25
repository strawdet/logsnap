package snapshot

import (
	"os"
	"testing"
	"time"
)

func writeRetentionSnapshot(t *testing.T, dir string, id string, createdAt time.Time) {
	t.Helper()
	s := &Snapshot{
		ID:        id,
		Label:     "retention-test",
		CreatedAt: createdAt,
		Entries:   sampleEntries(),
	}
	if err := s.Save(dir); err != nil {
		t.Fatalf("save snapshot %s: %v", id, err)
	}
}

func TestSetRetentionPolicy_AndGet(t *testing.T) {
	dir := t.TempDir()
	policy := RetentionPolicy{
		MaxCount:    10,
		MaxAgeDays:  30,
		ProtectPins: true,
	}
	if err := SetRetentionPolicy(dir, policy); err != nil {
		t.Fatalf("set: %v", err)
	}
	got, err := GetRetentionPolicy(dir)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.MaxCount != 10 || got.MaxAgeDays != 30 || !got.ProtectPins {
		t.Errorf("unexpected policy: %+v", got)
	}
	if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
		t.Error("timestamps should be set")
	}
}

func TestGetRetentionPolicy_NoFile(t *testing.T) {
	dir := t.TempDir()
	policy, err := GetRetentionPolicy(dir)
	if err != nil {
		t.Fatalf("expected no error for missing policy, got: %v", err)
	}
	if policy.MaxCount != 0 || policy.MaxAgeDays != 0 {
		t.Errorf("expected zero policy, got: %+v", policy)
	}
}

func TestRemoveRetentionPolicy(t *testing.T) {
	dir := t.TempDir()
	_ = SetRetentionPolicy(dir, RetentionPolicy{MaxCount: 5})
	if err := RemoveRetentionPolicy(dir); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if _, err := os.Stat(retentionPath(dir)); !os.IsNotExist(err) {
		t.Error("policy file should be deleted")
	}
}

func TestRemoveRetentionPolicy_NotFound(t *testing.T) {
	dir := t.TempDir()
	if err := RemoveRetentionPolicy(dir); err == nil {
		t.Error("expected error when no policy exists")
	}
}

func TestApplyRetentionPolicy_MaxCount(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().UTC()
	writeRetentionSnapshot(t, dir, "snap-001", now.Add(-1*time.Hour))
	writeRetentionSnapshot(t, dir, "snap-002", now.Add(-2*time.Hour))
	writeRetentionSnapshot(t, dir, "snap-003", now.Add(-3*time.Hour))

	_ = SetRetentionPolicy(dir, RetentionPolicy{MaxCount: 1})
	deleted, err := ApplyRetentionPolicy(dir, false)
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(deleted) != 2 {
		t.Errorf("expected 2 deleted, got %d", len(deleted))
	}
}

func TestApplyRetentionPolicy_DryRun(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().UTC()
	writeRetentionSnapshot(t, dir, "snap-a", now.Add(-10*24*time.Hour))
	writeRetentionSnapshot(t, dir, "snap-b", now.Add(-20*24*time.Hour))

	_ = SetRetentionPolicy(dir, RetentionPolicy{MaxAgeDays: 5})
	deleted, err := ApplyRetentionPolicy(dir, true)
	if err != nil {
		t.Fatalf("dry run: %v", err)
	}
	if len(deleted) != 2 {
		t.Errorf("expected 2 candidates, got %d", len(deleted))
	}
	// files should still exist after dry run
	snaps, _ := ListSnapshots(dir)
	if len(snaps) != 2 {
		t.Errorf("dry run should not delete files, got %d snapshots", len(snaps))
	}
}
