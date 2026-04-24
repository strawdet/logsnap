package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// TimelineEntry represents a single event in a snapshot's timeline.
type TimelineEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Detail    string    `json:"detail,omitempty"`
}

// Timeline holds an ordered list of events for a snapshot.
type Timeline struct {
	SnapshotID string          `json:"snapshot_id"`
	Entries    []TimelineEntry `json:"entries"`
}

func timelinePath(dir, snapshotID string) string {
	return filepath.Join(dir, snapshotID+".timeline.json")
}

// AddTimelineEvent appends an event to the snapshot's timeline.
func AddTimelineEvent(dir, snapshotID, event, detail string) error {
	snapshotFile := filepath.Join(dir, snapshotID+".json")
	if _, err := os.Stat(snapshotFile); os.IsNotExist(err) {
		return fmt.Errorf("snapshot %q not found", snapshotID)
	}

	tl, err := GetTimeline(dir, snapshotID)
	if err != nil {
		return err
	}

	tl.Entries = append(tl.Entries, TimelineEntry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Detail:    detail,
	})

	return saveTimeline(dir, tl)
}

// GetTimeline loads the timeline for a snapshot, returning an empty one if not yet created.
func GetTimeline(dir, snapshotID string) (*Timeline, error) {
	path := timelinePath(dir, snapshotID)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Timeline{SnapshotID: snapshotID, Entries: []TimelineEntry{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading timeline: %w", err)
	}
	var tl Timeline
	if err := json.Unmarshal(data, &tl); err != nil {
		return nil, fmt.Errorf("parsing timeline: %w", err)
	}
	sort.Slice(tl.Entries, func(i, j int) bool {
		return tl.Entries[i].Timestamp.Before(tl.Entries[j].Timestamp)
	})
	return &tl, nil
}

// ClearTimeline removes the timeline file for a snapshot.
func ClearTimeline(dir, snapshotID string) error {
	path := timelinePath(dir, snapshotID)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing timeline: %w", err)
	}
	return nil
}

func saveTimeline(dir string, tl *Timeline) error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding timeline: %w", err)
	}
	return os.WriteFile(timelinePath(dir, tl.SnapshotID), data, 0644)
}
