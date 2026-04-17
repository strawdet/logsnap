package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Highlight struct {
	SnapshotID string   `json:"snapshot_id"`
	Messages   []string `json:"messages"`
}

func highlightPath(dir, snapshotID string) string {
	return filepath.Join(dir, snapshotID+".highlight.json")
}

func AddHighlight(dir, snapshotID, message string) error {
	snap, err := Load(dir, snapshotID)
	if err != nil {
		return fmt.Errorf("snapshot not found: %w", err)
	}

	h := &Highlight{SnapshotID: snap.ID}
	p := highlightPath(dir, snapshotID)

	if data, err := os.ReadFile(p); err == nil {
		_ = json.Unmarshal(data, h)
	}

	for _, m := range h.Messages {
		if m == message {
			return nil
		}
	}
	h.Messages = append(h.Messages, message)

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

func GetHighlights(dir, snapshotID string) (*Highlight, error) {
	p := highlightPath(dir, snapshotID)
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return &Highlight{SnapshotID: snapshotID, Messages: []string{}}, nil
		}
		return nil, err
	}
	var h Highlight
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

func RemoveHighlight(dir, snapshotID, message string) error {
	h, err := GetHighlights(dir, snapshotID)
	if err != nil {
		return err
	}
	filtered := h.Messages[:0]
	for _, m := range h.Messages {
		if m != message {
			filtered = append(filtered, m)
		}
	}
	if len(filtered) == len(h.Messages) {
		return fmt.Errorf("highlight not found: %s", message)
	}
	h.Messages = filtered
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(highlightPath(dir, snapshotID), data, 0644)
}
