package snapshot

import (
	"os"
	"testing"
)

func writeLockSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	snap := &Snapshot{ID: id, Label: "lock-test"}
	snap.Entries = []LogEntry{{Level: "info", Message: "hello"}}
	p := dir + "/" + id + ".json"
	if err := writeJSON(p, snap); err != nil {
		t.Fatalf("writeJSON: %v", err)
	}
}

func TestLockSnapshot_AndIsLocked(t *testing.T) {
	dir := t.TempDir()
	writeLockSnapshot(t, dir, "snap1")

	if err := LockSnapshot(dir, "snap1", "do not delete"); err != nil {
		t.Fatalf("LockSnapshot: %v", err)
	}
	if !IsLocked(dir, "snap1") {
		t.Error("expected snap1 to be locked")
	}
}

func TestLockSnapshot_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := LockSnapshot(dir, "ghost", "")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestUnlockSnapshot(t *testing.T) {
	dir := t.TempDir()
	writeLockSnapshot(t, dir, "snap2")
	_ = LockSnapshot(dir, "snap2", "temp lock")

	if err := UnlockSnapshot(dir, "snap2"); err != nil {
		t.Fatalf("UnlockSnapshot: %v", err)
	}
	if IsLocked(dir, "snap2") {
		t.Error("expected snap2 to be unlocked")
	}
}

func TestUnlockSnapshot_NotLocked(t *testing.T) {
	dir := t.TempDir()
	writeLockSnapshot(t, dir, "snap3")
	err := UnlockSnapshot(dir, "snap3")
	if err == nil {
		t.Error("expected error when unlocking non-locked snapshot")
	}
}

func TestGetLockInfo(t *testing.T) {
	dir := t.TempDir()
	writeLockSnapshot(t, dir, "snap4")
	_ = LockSnapshot(dir, "snap4", "audit")

	info, err := GetLockInfo(dir, "snap4")
	if err != nil {
		t.Fatalf("GetLockInfo: %v", err)
	}
	if info.SnapshotID != "snap4" {
		t.Errorf("expected snap4, got %s", info.SnapshotID)
	}
	if info.Reason != "audit" {
		t.Errorf("expected reason 'audit', got %s", info.Reason)
	}
	if info.LockedAt.IsZero() {
		t.Error("expected non-zero LockedAt")
	}
}

func TestGetLockInfo_NotLocked(t *testing.T) {
	dir := t.TempDir()
	writeLockSnapshot(t, dir, "snap5")
	_, err := GetLockInfo(dir, "snap5")
	if err == nil {
		t.Error("expected error for unlocked snapshot")
	}
	_ = os.Remove(dir) // cleanup hint
}
