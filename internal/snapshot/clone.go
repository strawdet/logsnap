package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CloneSnapshot creates a copy of an existing snapshot with a new ID and optional label.
func CloneSnapshot(dir, id, newLabel string) (*Snapshot, error) {
	src, err := Load(dir, id)
	if err != nil {
		return nil, fmt.Errorf("clone: source snapshot not found: %w", err)
	}

	cloned := &Snapshot{
		ID:        generateID(),
		CreatedAt: time.Now().UTC(),
		Entries:   make([]LogEntry, len(src.Entries)),
	}

	copy(cloned.Entries, src.Entries)

	if newLabel != "" {
		cloned.Label = newLabel
	} else {
		cloned.Label = src.Label + " (clone)"
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("clone: failed to create dir: %w", err)
	}

	path := filepath.Join(dir, cloned.ID+".json")
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("clone: failed to create file: %w", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(cloned); err != nil {
		return nil, fmt.Errorf("clone: failed to write snapshot: %w", err)
	}

	return cloned, nil
}
