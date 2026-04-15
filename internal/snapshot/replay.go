package snapshot

import (
	"fmt"
	"io"
	"os"
	"time"
)

// ReplayOptions controls how a snapshot is replayed.
type ReplayOptions struct {
	Delay    time.Duration // delay between entries
	Filter   string        // optional level filter (e.g. "error")
	Writer   io.Writer     // output destination, defaults to os.Stdout
}

// ReplayResult holds the outcome of a replay operation.
type ReplayResult struct {
	Replayed int
	Skipped  int
}

// Replay streams log entries from a snapshot to the given writer,
// optionally filtering by level and inserting a delay between entries.
func Replay(snap *Snapshot, opts ReplayOptions) (ReplayResult, error) {
	if snap == nil {
		return ReplayResult{}, fmt.Errorf("snapshot is nil")
	}

	out := opts.Writer
	if out == nil {
		out = os.Stdout
	}

	var result ReplayResult

	for _, entry := range snap.Entries {
		if opts.Filter != "" && entry.Level != opts.Filter {
			result.Skipped++
			continue
		}

		line := fmt.Sprintf("[%s] %s: %s\n", entry.Timestamp.Format(time.RFC3339), entry.Level, entry.Message)
		if _, err := fmt.Fprint(out, line); err != nil {
			return result, fmt.Errorf("write error: %w", err)
		}

		result.Replayed++

		if opts.Delay > 0 {
			time.Sleep(opts.Delay)
		}
	}

	return result, nil
}
