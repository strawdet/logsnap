package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const baselineIndexFile = "baseline_index.json"

// BaselineIndex maps a named baseline to a snapshot ID.
type BaselineIndex map[string]string

func baselineIndexPath(dir string) string {
	return filepath.Join(dir, baselineIndexFile)
}

// LoadBaselineIndex reads the baseline index from disk.
func LoadBaselineIndex(dir string) (BaselineIndex, error) {
	path := baselineIndexPath(dir)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return BaselineIndex{}, nil
	}
	if err != nil {
		return nil, err
	}
	var idx BaselineIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return idx, nil
}

// SaveBaselineIndex writes the baseline index to disk.
func SaveBaselineIndex(dir string, idx BaselineIndex) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(baselineIndexPath(dir), data, 0644)
}

// SetBaseline assigns a snapshot ID to a named baseline.
func SetBaseline(dir, name, snapshotID string) error {
	path := filepath.Join(dir, snapshotID+".json")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return errors.New("snapshot not found: " + snapshotID)
	}
	idx, err := LoadBaselineIndex(dir)
	if err != nil {
		return err
	}
	idx[name] = snapshotID
	return SaveBaselineIndex(dir, idx)
}

// ResolveBaseline returns the snapshot ID for a named baseline.
func ResolveBaseline(dir, name string) (string, error) {
	idx, err := LoadBaselineIndex(dir)
	if err != nil {
		return "", err
	}
	id, ok := idx[name]
	if !ok {
		return "", errors.New("baseline not found: " + name)
	}
	return id, nil
}

// RemoveBaseline deletes a named baseline from the index.
func RemoveBaseline(dir, name string) error {
	idx, err := LoadBaselineIndex(dir)
	if err != nil {
		return err
	}
	if _, ok := idx[name]; !ok {
		return errors.New("baseline not found: " + name)
	}
	delete(idx, name)
	return SaveBaselineIndex(dir, idx)
}
