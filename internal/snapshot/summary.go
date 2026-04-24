package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// Summary holds a high-level overview of a snapshot.
type Summary struct {
	ID        string            `json:"id"`
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Total     int               `json:"total"`
	BySeverity map[string]int   `json:"by_severity"`
	TopMessages []string        `json:"top_messages"`
	Tags      []string          `json:"tags"`
}

// SummarizeSnapshot produces a Summary from a Snapshot.
func SummarizeSnapshot(s *Snapshot, topN int) (*Summary, error) {
	if s == nil {
		return nil, fmt.Errorf("snapshot is nil")
	}

	bySeverity := make(map[string]int)
	msgCount := make(map[string]int)

	for _, e := range s.Entries {
		if e.Level != "" {
			bySeverity[e.Level]++
		}
		if e.Message != "" {
			msgCount[e.Message]++
		}
	}

	topMessages := topNMessages(msgCount, topN)

	tags, _ := GetSnapshotTags(s.Dir, s.ID)

	return &Summary{
		ID:          s.ID,
		Label:       s.Label,
		CreatedAt:   s.CreatedAt,
		Total:       len(s.Entries),
		BySeverity:  bySeverity,
		TopMessages: topMessages,
		Tags:        tags,
	}, nil
}

// topNMessages returns the top n most frequent messages.
func topNMessages(counts map[string]int, n int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range counts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	var result []string
	for i, item := range sorted {
		if i >= n {
			break
		}
		result = append(result, item.Key)
	}
	return result
}

// GetSnapshotTags returns tags associated with a snapshot ID, or nil if none.
func GetSnapshotTags(dir, id string) ([]string, error) {
	index, err := LoadTagIndex(dir)
	if err != nil {
		return nil, err
	}
	var tags []string
	for tag, sid := range index {
		if sid == id {
			tags = append(tags, tag)
		}
	}
	sort.Strings(tags)
	return tags, nil
}
