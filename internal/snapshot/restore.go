package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// RestoreResult holds the outcome of a restore operation.
type RestoreResult struct {
	SnapshotID string
	OutputPath string
	EntryCount int
}

// Restore writes the log entries from a snapshot back to a log file.
// If outputPath is "-", entries are written to stdout as newline-delimited JSON.
func Restore(dir, snapshotID, outputPath string) (*RestoreResult, error) {
	snap, err := Load(dir, snapshotID)
	if err != nil {
		return nil, fmt.Errorf("restore: load snapshot: %w", err)
	}

	var out *os.File
	if outputPath == "-" {
		out = os.Stdout
	} else {
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return nil, fmt.Errorf("restore: mkdir: %w", err)
		}
		f, err := os.Create(outputPath)
		if err != nil {
			return nil, fmt.Errorf("restore: create file: %w", err)
		}
		defer f.Close()
		out = f
	}

	enc := json.NewEncoder(out)
	for _, entry := range snap.Entries {
		if err := enc.Encode(entry); err != nil {
			return nil, fmt.Errorf("restore: encode entry: %w", err)
		}
	}

	return &RestoreResult{
		SnapshotID: snap.ID,
		OutputPath: outputPath,
		EntryCount: len(snap.Entries),
	}, nil
}
