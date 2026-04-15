package diff

import (
	"fmt"
	"strings"

	"github.com/user/logsnap/internal/snapshot"
)

// Result holds the outcome of comparing two snapshots.
type Result struct {
	Added   []snapshot.LogEntry
	Removed []snapshot.LogEntry
	Changed []Change
}

// Change represents a log entry whose level changed between snapshots.
type Change struct {
	Message  string
	OldLevel string
	NewLevel string
}

// Compare computes the diff between a baseline and a current snapshot.
func Compare(baseline, current *snapshot.Snapshot) *Result {
	result := &Result{}

	baselineMap := indexByMessage(baseline.Entries)
	currentMap := indexByMessage(current.Entries)

	for msg, cur := range currentMap {
		if base, found := baselineMap[msg]; !found {
			result.Added = append(result.Added, cur)
		} else if base.Level != cur.Level {
			result.Changed = append(result.Changed, Change{
				Message:  msg,
				OldLevel: base.Level,
				NewLevel: cur.Level,
			})
		}
	}

	for msg, base := range baselineMap {
		if _, found := currentMap[msg]; !found {
			result.Removed = append(result.Removed, base)
		}
	}

	return result
}

// indexByMessage builds a map of log entries keyed by their message.
func indexByMessage(entries []snapshot.LogEntry) map[string]snapshot.LogEntry {
	m := make(map[string]snapshot.LogEntry, len(entries))
	for _, e := range entries {
		m[e.Message] = e
	}
	return m
}

// Summary returns a human-readable summary of the diff result.
func (r *Result) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Added: %d  Removed: %d  Changed: %d\n", len(r.Added), len(r.Removed), len(r.Changed))

	for _, e := range r.Added {
		fmt.Fprintf(&sb, "  [+] [%s] %s\n", e.Level, e.Message)
	}
	for _, e := range r.Removed {
		fmt.Fprintf(&sb, "  [-] [%s] %s\n", e.Level, e.Message)
	}
	for _, c := range r.Changed {
		fmt.Fprintf(&sb, "  [~] %s  (%s -> %s)\n", c.Message, c.OldLevel, c.NewLevel)
	}

	return sb.String()
}
