package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RetentionPolicy defines rules for automatic snapshot cleanup.
type RetentionPolicy struct {
	MaxCount    int           `json:"max_count,omitempty"`    // keep at most N snapshots
	MaxAgeDays  int           `json:"max_age_days,omitempty"` // delete snapshots older than N days
	ProtectPins bool          `json:"protect_pins"`           // skip pinned snapshots
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func retentionPath(dir string) string {
	return filepath.Join(dir, "retention_policy.json")
}

// SetRetentionPolicy saves a retention policy to the snapshot directory.
func SetRetentionPolicy(dir string, policy RetentionPolicy) error {
	now := time.Now().UTC()
	if policy.CreatedAt.IsZero() {
		policy.CreatedAt = now
	}
	policy.UpdatedAt = now

	data, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return fmt.Errorf("retention: marshal: %w", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("retention: mkdir: %w", err)
	}
	return os.WriteFile(retentionPath(dir), data, 0644)
}

// GetRetentionPolicy loads the retention policy from the snapshot directory.
// Returns a zero-value policy and no error if no policy is set.
func GetRetentionPolicy(dir string) (RetentionPolicy, error) {
	var policy RetentionPolicy
	data, err := os.ReadFile(retentionPath(dir))
	if os.IsNotExist(err) {
		return policy, nil
	}
	if err != nil {
		return policy, fmt.Errorf("retention: read: %w", err)
	}
	if err := json.Unmarshal(data, &policy); err != nil {
		return policy, fmt.Errorf("retention: unmarshal: %w", err)
	}
	return policy, nil
}

// RemoveRetentionPolicy deletes the retention policy file.
func RemoveRetentionPolicy(dir string) error {
	err := os.Remove(retentionPath(dir))
	if os.IsNotExist(err) {
		return fmt.Errorf("retention: no policy set")
	}
	return err
}

// ApplyRetentionPolicy enforces the policy against the snapshot directory.
// Returns the list of snapshot IDs that were deleted.
func ApplyRetentionPolicy(dir string, dryRun bool) ([]string, error) {
	policy, err := GetRetentionPolicy(dir)
	if err != nil {
		return nil, err
	}

	snapshots, err := ListSnapshots(dir)
	if err != nil {
		return nil, err
	}

	pinIndex, _ := LoadPinIndex(dir)

	var toDelete []string
	cutoff := time.Now().UTC().AddDate(0, 0, -policy.MaxAgeDays)

	for i, s := range snapshots {
		if policy.ProtectPins {
			if _, pinned := pinIndex[s.ID]; pinned {
				continue
			}
		}
		if policy.MaxAgeDays > 0 && s.CreatedAt.Before(cutoff) {
			toDelete = append(toDelete, s.ID)
			continue
		}
		if policy.MaxCount > 0 && i >= policy.MaxCount {
			toDelete = append(toDelete, s.ID)
		}
	}

	if dryRun {
		return toDelete, nil
	}
	for _, id := range toDelete {
		if err := DeleteSnapshot(dir, id); err != nil {
			return toDelete, fmt.Errorf("retention: delete %s: %w", id, err)
		}
	}
	return toDelete, nil
}
