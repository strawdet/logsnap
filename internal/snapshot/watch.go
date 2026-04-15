package snapshot

import (
	"context"
	"encoding/json"
	"os"
	"time"
)

// WatchOptions configures the watch behavior.
type WatchOptions struct {
	LogFile  string
	Dir      string
	Interval time.Duration
	Labels   map[string]string
	OnSnap   func(id string, err error)
}

// Watch polls a log file at the given interval, capturing a new snapshot
// each time the file changes. It blocks until ctx is cancelled.
func Watch(ctx context.Context, opts WatchOptions) error {
	if opts.Interval <= 0 {
		opts.Interval = 30 * time.Second
	}

	var lastSize int64
	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			info, err := os.Stat(opts.LogFile)
			if err != nil {
				if opts.OnSnap != nil {
					opts.OnSnap("", err)
				}
				continue
			}

			if info.Size() == lastSize {
				continue
			}
			lastSize = info.Size()

			entries, err := readEntries(opts.LogFile)
			if err != nil {
				if opts.OnSnap != nil {
					opts.OnSnap("", err)
				}
				continue
			}

			snap := New(entries, opts.Labels)
			if err := snap.Save(opts.Dir); err != nil {
				if opts.OnSnap != nil {
					opts.OnSnap("", err)
				}
				continue
			}

			if opts.OnSnap != nil {
				opts.OnSnap(snap.ID, nil)
			}
		}
	}
}

// readEntries reads newline-delimited JSON log entries from path.
func readEntries(path string) ([]LogEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []LogEntry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e LogEntry
		if err := dec.Decode(&e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries, nil
}
