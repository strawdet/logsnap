package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type LockInfo struct {
	SnapshotID string    `json:"snapshot_id"`
	LockedAt   time.Time `json:"locked_at"`
	Reason     string    `json:"reason,omitempty"`
}

func lockPath(dir, id string) string {
	return filepath.Join(dir, id+".lock")
}

// LockSnapshot marks a snapshot as locked, preventing deletion or modification.
func LockSnapshot(dir, id, reason string) error {
	snapshotFile := filepath.Join(dir, id+".json")
	if _, err := os.Stat(snapshotFile); os.IsNotExist(err) {
		return fmt.Errorf("snapshot %q not found", id)
	}

	info := LockInfo{
		SnapshotID: id,
		LockedAt:   time.Now().UTC(),
		Reason:     reason,
	}
	return writeJSON(lockPath(dir, id), info)
}

// UnlockSnapshot removes the lock from a snapshot.
func UnlockSnapshot(dir, id string) error {
	p := lockPath(dir, id)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("snapshot %q is not locked", id)
	}
	return os.Remove(p)
}

// IsLocked returns true if the snapshot has a lock file.
func IsLocked(dir, id string) bool {
	_, err := os.Stat(lockPath(dir, id))
	return err == nil
}

// GetLockInfo returns the LockInfo for a locked snapshot.
func GetLockInfo(dir, id string) (*LockInfo, error) {
	p := lockPath(dir, id)
	var info LockInfo
	if err := readJSON(p, &info); err != nil {
		return nil, fmt.Errorf("snapshot %q is not locked or lock unreadable: %w", id, err)
	}
	return &info, nil
}
