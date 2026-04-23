package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type LabelIndex struct {
	Labels map[string][]string `json:"labels"` // label -> []snapshotID
}

func labelIndexPath(dir string) string {
	return filepath.Join(dir, ".label_index.json")
}

func LoadLabelIndex(dir string) (*LabelIndex, error) {
	path := labelIndexPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &LabelIndex{Labels: make(map[string][]string)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read label index: %w", err)
	}
	var idx LabelIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse label index: %w", err)
	}
	if idx.Labels == nil {
		idx.Labels = make(map[string][]string)
	}
	return &idx, nil
}

func SaveLabelIndex(dir string, idx *LabelIndex) error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal label index: %w", err)
	}
	return os.WriteFile(labelIndexPath(dir), data, 0644)
}

func AddLabel(dir, snapshotID, label string) error {
	snap, err := Load(dir, snapshotID)
	if err != nil {
		return fmt.Errorf("snapshot not found: %w", err)
	}
	_ = snap

	idx, err := LoadLabelIndex(dir)
	if err != nil {
		return err
	}
	for _, id := range idx.Labels[label] {
		if id == snapshotID {
			return nil // already labeled
		}
	}
	idx.Labels[label] = append(idx.Labels[label], snapshotID)
	return SaveLabelIndex(dir, idx)
}

func RemoveLabel(dir, snapshotID, label string) error {
	idx, err := LoadLabelIndex(dir)
	if err != nil {
		return err
	}
	ids := idx.Labels[label]
	updated := ids[:0]
	for _, id := range ids {
		if id != snapshotID {
			updated = append(updated, id)
		}
	}
	if len(updated) == 0 {
		delete(idx.Labels, label)
	} else {
		idx.Labels[label] = updated
	}
	return SaveLabelIndex(dir, idx)
}

func GetSnapshotsByLabel(dir, label string) ([]string, error) {
	idx, err := LoadLabelIndex(dir)
	if err != nil {
		return nil, err
	}
	return idx.Labels[label], nil
}

func ListLabelsForSnapshot(dir, snapshotID string) ([]string, error) {
	idx, err := LoadLabelIndex(dir)
	if err != nil {
		return nil, err
	}
	var labels []string
	for label, ids := range idx.Labels {
		for _, id := range ids {
			if id == snapshotID {
				labels = append(labels, label)
				break
			}
		}
	}
	return labels, nil
}
