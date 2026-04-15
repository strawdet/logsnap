package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Meta holds lightweight metadata about a stored snapshot.
type Meta struct {
	ID        string    `json:"id"`
	Label     string    `json:"label"`
	CreatedAt time.Time `json:"created_at"`
	Entries   int       `json:"entries"`
}

// ListSnapshots returns metadata for all snapshots found in dir,
// sorted by creation time (newest first).
func ListSnapshots(dir string) ([]Meta, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading snapshot dir: %w", err)
	}

	var metas []Meta
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}

		path := filepath.Join(dir, e.Name())
		snap, err := Load(path)
		if err != nil {
			continue // skip unreadable files
		}

		metas = append(metas, Meta{
			ID:        snap.ID,
			Label:     snap.Label,
			CreatedAt: snap.CreatedAt,
			Entries:   len(snap.Entries),
		})
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].CreatedAt.After(metas[j].CreatedAt)
	})

	return metas, nil
}

// DeleteSnapshot removes a snapshot file by ID from dir.
func DeleteSnapshot(dir, id string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading snapshot dir: %w", err)
	}

	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var snap struct {
			ID string `json:"id"`
		}{}
		if json.Unmarshal(data, &snap) == nil && snap.ID == id {
			return os.Remove(path)
		}
	}
	return fmt.Errorf("snapshot %q not found", id)
}
