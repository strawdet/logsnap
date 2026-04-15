package snapshot

import (
	"fmt"
	"sort"
)

// LevelCounts holds the count of log entries per level.
type LevelCounts map[string]int

// Stats holds aggregate statistics for a snapshot.
type Stats struct {
	SnapshotID  string
	Label       string
	TotalCount  int
	LevelCounts LevelCounts
	TopMessages []string
}

// ComputeStats returns statistics derived from the given snapshot.
func ComputeStats(snap *Snapshot) (*Stats, error) {
	if snap == nil {
		return nil, fmt.Errorf("snapshot must not be nil")
	}

	levelCounts := make(LevelCounts)
	msgFreq := make(map[string]int)

	for _, entry := range snap.Entries {
		level := entry.Level
		if level == "" {
			level = "unknown"
		}
		levelCounts[level]++
		msgFreq[entry.Message]++
	}

	top := topN(msgFreq, 5)

	return &Stats{
		SnapshotID:  snap.ID,
		Label:       snap.Label,
		TotalCount:  len(snap.Entries),
		LevelCounts: levelCounts,
		TopMessages: top,
	}, nil
}

// topN returns up to n messages sorted by frequency descending.
func topN(freq map[string]int, n int) []string {
	type pair struct {
		msg   string
		count int
	}
	pairs := make([]pair, 0, len(freq))
	for msg, count := range freq {
		pairs = append(pairs, pair{msg, count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].msg < pairs[j].msg
	})
	result := make([]string, 0, n)
	for i := 0; i < n && i < len(pairs); i++ {
		result = append(result, pairs[i].msg)
	}
	return result
}
