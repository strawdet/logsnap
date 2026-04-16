package snapshot

import (
	"strings"
	"time"
)

// SearchFilter defines criteria for filtering snapshots.
type SearchFilter struct {
	Tag       string
	Since     *time.Time
	Until     *time.Time
	LabelKey  string
	LabelVal  string
}

// SearchResult holds a matched snapshot and its metadata.
type SearchResult struct {
	Snapshot *Snapshot
	FilePath string
}

// Search returns snapshots from dir that match the given filter.
func Search(dir string, filter SearchFilter) ([]SearchResult, error) {
	snaps, err := ListSnapshots(dir)
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, meta := range snaps {
		snap, err := Load(meta.FilePath)
		if err != nil {
			continue
		}

		if !matchesFilter(snap, filter) {
			continue
		}

		results = append(results, SearchResult{Snapshot: snap, FilePath: meta.FilePath})
	}

	return results, nil
}

// matchesFilter reports whether a snapshot satisfies all criteria in filter.
func matchesFilter(snap *Snapshot, filter SearchFilter) bool {
	if filter.Tag != "" {
		tagged := false
		for _, t := range snap.Tags {
			if strings.EqualFold(t, filter.Tag) {
				tagged = true
				break
			}
		}
		if !tagged {
			return false
		}
	}

	if filter.Since != nil && snap.CreatedAt.Before(*filter.Since) {
		return false
	}

	if filter.Until != nil && snap.CreatedAt.After(*filter.Until) {
		return false
	}

	if filter.LabelKey != "" {
		val, ok := snap.Labels[filter.LabelKey]
		if !ok {
			return false
		}
		if filter.LabelVal != "" && !strings.EqualFold(val, filter.LabelVal) {
			return false
		}
	}

	return true
}
