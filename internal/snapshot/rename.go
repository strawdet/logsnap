package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// RenameSnapshot renames a snapshot by updating its label in the metadata file.
func RenameSnapshot(dir, id, newLabel string) error {
	path := filepath.Join(dir, id+".json")

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("snapshot %q not found", id)
		}
		return fmt.Errorf("read snapshot: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return fmt.Errorf("parse snapshot: %w", err)
	}

	snap.Label = newLabel

	updated, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	if err := os.WriteFile(path, updated, 0644); err != nil {
		return fmt.Errorf("write snapshot: %w", err)
	}

	return nil
}
