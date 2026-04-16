package snapshot

import (
	"fmt"
	"path/filepath"
	"time"
)

// MergeResult holds the outcome of merging two snapshots.
type MergeResult struct {
	Merged    *Snapshot
	Conflicts []string
}

// MergeSnapshots combines entries from two snapshots into a new one.
// Entries with duplicate messages are flagged as conflicts; the base entry is kept.
func MergeSnapshots(dir, baseID, otherID, label string) (*MergeResult, error) {
	base, err := Load(dir, baseID)
	if err != nil {
		return nil, fmt.Errorf("loading base snapshot: %w", err)
	}
	other, err := Load(dir, otherID)
	if err != nil {
		return nil, fmt.Errorf("loading other snapshot: %w", err)
	}

	seen := make(map[string]bool)
	var merged []LogEntry
	var conflicts []string

	for _, e := range base.Entries {
		seen[e.Message] = true
		merged = append(merged, e)
	}

	for _, e := range other.Entries {
		if seen[e.Message] {
			conflicts = append(conflicts, e.Message)
			continue
		}
		seen[e.Message] = true
		merged = append(merged, e)
	}

	if label == "" {
		label = fmt.Sprintf("merge-%s-%s", baseID[:8], otherID[:8])
	}

	snap := &Snapshot{
		ID:        generateID(),
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Entries:   merged,
	}

	path := filepath.Join(dir, snap.ID+".json")
	if err := snap.Save(path); err != nil {
		return nil, fmt.Errorf("saving merged snapshot: %w", err)
	}

	return &MergeResult{Merged: snap, Conflicts: conflicts}, nil
}
