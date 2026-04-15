package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CompareResult holds a named comparison between two snapshots.
type CompareResult struct {
	BaseID   string `json:"base_id"`
	TargetID string `json:"target_id"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}

// SaveCompareResult persists a named comparison result to disk.
func SaveCompareResult(dir, name, baseID, targetID string) (*CompareResult, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create compare dir: %w", err)
	}

	result := &CompareResult{
		BaseID:   baseID,
		TargetID: targetID,
		Name:     name,
		FilePath: filepath.Join(dir, name+".json"),
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal compare result: %w", err)
	}

	if err := os.WriteFile(result.FilePath, data, 0644); err != nil {
		return nil, fmt.Errorf("write compare result: %w", err)
	}

	return result, nil
}

// LoadCompareResult loads a named comparison result from disk.
func LoadCompareResult(dir, name string) (*CompareResult, error) {
	path := filepath.Join(dir, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("compare result %q not found", name)
		}
		return nil, fmt.Errorf("read compare result: %w", err)
	}

	var result CompareResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse compare result: %w", err)
	}

	return &result, nil
}

// ListCompareResults returns all saved comparison names in the given directory.
func ListCompareResults(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("read compare dir: %w", err)
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
