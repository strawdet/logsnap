package snapshot

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeExportSnapshot() *Snapshot {
	return &Snapshot{
		ID:        "export-test-id",
		CreatedAt: time.Now(),
		Entries: []LogEntry{
			{Timestamp: time.Now(), Level: "info", Message: "service started", Fields: map[string]interface{}{"pid": 42}},
			{Timestamp: time.Now(), Level: "error", Message: "connection refused", Fields: map[string]interface{}{"host": "db"}},
		},
	}
}

func TestExportSnapshot_JSON(t *testing.T) {
	snap := makeExportSnapshot()
	dest := filepath.Join(t.TempDir(), "out.json")

	if err := ExportSnapshot(snap, dest, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("reading output file: %v", err)
	}

	var entries []LogEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("invalid json output: %v", err)
	}
	if len(entries) != len(snap.Entries) {
		t.Errorf("expected %d entries, got %d", len(snap.Entries), len(entries))
	}
}

func TestExportSnapshot_CSV(t *testing.T) {
	snap := makeExportSnapshot()
	dest := filepath.Join(t.TempDir(), "out.csv")

	if err := ExportSnapshot(snap, dest, FormatCSV); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f, err := os.Open(dest)
	if err != nil {
		t.Fatalf("opening csv: %v", err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("reading csv: %v", err)
	}
	// header + data rows
	if len(rows) != len(snap.Entries)+1 {
		t.Errorf("expected %d rows (incl. header), got %d", len(snap.Entries)+1, len(rows))
	}
	if rows[0][0] != "timestamp" {
		t.Errorf("expected header row, got %v", rows[0])
	}
}

func TestExportSnapshot_UnsupportedFormat(t *testing.T) {
	snap := makeExportSnapshot()
	dest := filepath.Join(t.TempDir(), "out.xml")

	err := ExportSnapshot(snap, dest, "xml")
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}
