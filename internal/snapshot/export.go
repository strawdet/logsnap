package snapshot

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ExportFormat represents the output format for snapshot export.
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
)

// ExportSnapshot writes a snapshot's log entries to the given file path
// in the specified format (json or csv).
func ExportSnapshot(snap *Snapshot, destPath string, format ExportFormat) error {
	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("creating export file: %w", err)
	}
	defer f.Close()

	switch format {
	case FormatJSON:
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(snap.Entries); err != nil {
			return fmt.Errorf("encoding json: %w", err)
		}
	case FormatCSV:
		w := csv.NewWriter(f)
		if err := w.Write([]string{"timestamp", "level", "message", "fields"}); err != nil {
			return fmt.Errorf("writing csv header: %w", err)
		}
		for _, e := range snap.Entries {
			fieldsJSON, _ := json.Marshal(e.Fields)
			row := []string{
				e.Timestamp.Format(time.RFC3339),
				e.Level,
				e.Message,
				string(fieldsJSON),
			}
			if err := w.Write(row); err != nil {
				return fmt.Errorf("writing csv row: %w", err)
			}
		}
		w.Flush()
		if err := w.Error(); err != nil {
			return fmt.Errorf("flushing csv: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format %q: use \"json\" or \"csv\"", format)
	}

	return nil
}
