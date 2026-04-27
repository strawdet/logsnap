package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// VersionEntry represents a single version record for a snapshot.
type VersionEntry struct {
	Version   int       `json:"version"`
	SnapshotID string   `json:"snapshot_id"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	Note      string    `json:"note,omitempty"`
}

// VersionIndex holds all version entries for a snapshot.
type VersionIndex struct {
	Entries []VersionEntry `json:"entries"`
}

func versionPath(dir, snapshotID string) string {
	return filepath.Join(dir, snapshotID+".versions.json")
}

// AddVersion records a new version entry for the given snapshot.
func AddVersion(dir, snapshotID, note string) error {
	snap, err := Load(dir, snapshotID)
	if err != nil {
		return fmt.Errorf("snapshot not found: %w", err)
	}

	index, _ := LoadVersionIndex(dir, snapshotID)
	nextVersion := len(index.Entries) + 1

	entry := VersionEntry{
		Version:    nextVersion,
		SnapshotID: snap.ID,
		Label:      snap.Label,
		CreatedAt:  time.Now().UTC(),
		Note:       note,
	}
	index.Entries = append(index.Entries, entry)

	return saveVersionIndex(dir, snapshotID, index)
}

// LoadVersionIndex loads the version history for a snapshot.
func LoadVersionIndex(dir, snapshotID string) (VersionIndex, error) {
	path := versionPath(dir, snapshotID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return VersionIndex{}, nil
		}
		return VersionIndex{}, err
	}
	var index VersionIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return VersionIndex{}, err
	}
	return index, nil
}

// ClearVersionHistory removes all version records for a snapshot.
func ClearVersionHistory(dir, snapshotID string) error {
	path := versionPath(dir, snapshotID)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func saveVersionIndex(dir, snapshotID string, index VersionIndex) error {
	path := versionPath(dir, snapshotID)
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
