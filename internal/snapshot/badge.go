package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Badge represents a named achievement or status marker attached to a snapshot.
type Badge struct {
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Reason    string    `json:"reason"`
	AwardedAt time.Time `json:"awarded_at"`
}

type BadgeIndex struct {
	Badges []Badge `json:"badges"`
}

func badgePath(dir, snapshotID string) string {
	return filepath.Join(dir, snapshotID+".badges.json")
}

// AddBadge attaches a badge to the given snapshot.
func AddBadge(dir, snapshotID, name, icon, reason string) error {
	snaps, err := ListSnapshots(dir)
	if err != nil {
		return err
	}
	found := false
	for _, s := range snaps {
		if s.ID == snapshotID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("snapshot %q not found", snapshotID)
	}

	index, _ := GetBadges(dir, snapshotID)
	for _, b := range index.Badges {
		if b.Name == name {
			return nil // already awarded
		}
	}
	index.Badges = append(index.Badges, Badge{
		Name:      name,
		Icon:      icon,
		Reason:    reason,
		AwardedAt: time.Now().UTC(),
	})
	return saveBadgeIndex(dir, snapshotID, index)
}

// GetBadges returns all badges for a snapshot.
func GetBadges(dir, snapshotID string) (BadgeIndex, error) {
	path := badgePath(dir, snapshotID)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return BadgeIndex{}, nil
	}
	if err != nil {
		return BadgeIndex{}, err
	}
	var index BadgeIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return BadgeIndex{}, err
	}
	return index, nil
}

// RemoveBadge removes a badge by name from a snapshot.
func RemoveBadge(dir, snapshotID, name string) error {
	index, err := GetBadges(dir, snapshotID)
	if err != nil {
		return err
	}
	filtered := index.Badges[:0]
	for _, b := range index.Badges {
		if b.Name != name {
			filtered = append(filtered, b)
		}
	}
	if len(filtered) == len(index.Badges) {
		return fmt.Errorf("badge %q not found on snapshot %q", name, snapshotID)
	}
	index.Badges = filtered
	return saveBadgeIndex(dir, snapshotID, index)
}

func saveBadgeIndex(dir, snapshotID string, index BadgeIndex) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(badgePath(dir, snapshotID), data, 0644)
}
