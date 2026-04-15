package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const compareIndexFile = "compare_index.json"

// CompareIndex maps comparison names to their metadata.
type CompareIndex map[string]CompareResult

// compareIndexPath returns the path to the compare index file.
func compareIndexPath(dir string) string {
	return filepath.Join(dir, compareIndexFile)
}

// LoadCompareIndex reads the compare index from disk.
func LoadCompareIndex(dir string) (CompareIndex, error) {
	path := compareIndexPath(dir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return make(CompareIndex), nil
		}
		return nil, fmt.Errorf("read compare index: %w", err)
	}

	var index CompareIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("parse compare index: %w", err)
	}
	return index, nil
}

// SaveCompareIndex writes the compare index to disk.
func SaveCompareIndex(dir string, index CompareIndex) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create compare dir: %w", err)
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal compare index: %w", err)
	}

	return os.WriteFile(compareIndexPath(dir), data, 0644)
}

// RegisterCompare adds or updates a comparison entry in the index.
func RegisterCompare(dir, name, baseID, targetID string) error {
	index, err := LoadCompareIndex(dir)
	if err != nil {
		return err
	}

	index[name] = CompareResult{
		Name:     name,
		BaseID:   baseID,
		TargetID: targetID,
		FilePath: filepath.Join(dir, name+".json"),
	}

	return SaveCompareIndex(dir, index)
}

// DeregisterCompare removes a comparison entry from the index.
func DeregisterCompare(dir, name string) error {
	index, err := LoadCompareIndex(dir)
	if err != nil {
		return err
	}

	if _, ok := index[name]; !ok {
		return fmt.Errorf("comparison %q not found in index", name)
	}

	delete(index, name)
	return SaveCompareIndex(dir, index)
}
