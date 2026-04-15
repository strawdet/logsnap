package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// TagIndex maps tag names to snapshot IDs.
type TagIndex map[string]string

func tagIndexPath(dir string) string {
	return filepath.Join(dir, "tags.json")
}

// LoadTagIndex reads the tag index from disk. Returns an empty index if not found.
func LoadTagIndex(dir string) (TagIndex, error) {
	path := tagIndexPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(TagIndex), nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading tag index: %w", err)
	}
	var index TagIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("parsing tag index: %w", err)
	}
	return index, nil
}

// SaveTagIndex writes the tag index to disk.
func SaveTagIndex(dir string, index TagIndex) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling tag index: %w", err)
	}
	return os.WriteFile(tagIndexPath(dir), data, 0644)
}

// TagSnapshot associates a human-readable tag with a snapshot ID.
func TagSnapshot(dir, tag, snapshotID string) error {
	index, err := LoadTagIndex(dir)
	if err != nil {
		return err
	}
	index[tag] = snapshotID
	return SaveTagIndex(dir, index)
}

// ResolveTag returns the snapshot ID for a given tag, or an error if not found.
func ResolveTag(dir, tag string) (string, error) {
	index, err := LoadTagIndex(dir)
	if err != nil {
		return "", err
	}
	id, ok := index[tag]
	if !ok {
		return "", fmt.Errorf("tag %q not found", tag)
	}
	return id, nil
}

// RemoveTag deletes a tag from the index.
func RemoveTag(dir, tag string) error {
	index, err := LoadTagIndex(dir)
	if err != nil {
		return err
	}
	if _, ok := index[tag]; !ok {
		return fmt.Errorf("tag %q not found", tag)
	}
	delete(index, tag)
	return SaveTagIndex(dir, index)
}
