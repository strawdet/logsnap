package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// PruneOptions controls which snapshots are removed during pruning.
type PruneOptions struct {
	// KeepLast retains the N most recent snapshots (0 = no limit).
	KeepLast int
	// OlderThan removes snapshots created before this time (zero = no limit).
	OlderThan time.Time
	// DryRun lists what would be deleted without actually deleting.
	DryRun bool
}

// PruneResult summarises what was (or would be) removed.
type PruneResult struct {
	Removed []string
	Kept    []string
}

// Prune removes snapshots from dir according to opts.
func Prune(dir string, opts PruneOptions) (*PruneResult, error) {
	snapshots, err := ListSnapshots(dir)
	if err != nil {
		return nil, fmt.Errorf("prune: list snapshots: %w", err)
	}

	// ListSnapshots returns newest-first; keep that order.
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].CreatedAt.After(snapshots[j].CreatedAt)
	})

	result := &PruneResult{}

	for i, s := range snapshots {
		should := false

		if opts.KeepLast > 0 && i >= opts.KeepLast {
			should = true
		}
		if !opts.OlderThan.IsZero() && s.CreatedAt.Before(opts.OlderThan) {
			should = true
		}

		if should {
			result.Removed = append(result.Removed, s.ID)
			if !opts.DryRun {
				path := filepath.Join(dir, s.ID+".json")
				if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
					return result, fmt.Errorf("prune: remove %s: %w", s.ID, err)
				}
			}
		} else {
			result.Kept = append(result.Kept, s.ID)
		}
	}

	return result, nil
}
