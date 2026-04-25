package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func writeAccessSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	f := filepath.Join(dir, id+".json")
	if err := os.WriteFile(f, []byte(`{"id":"`+id+`","entries":[]}`), 0644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}
}

func TestRecordAccess_AndGet(t *testing.T) {
	dir := t.TempDir()
	writeAccessSnapshot(t, dir, "snap1")

	if err := RecordAccess(dir, "snap1", "read", "alice", "initial read"); err != nil {
		t.Fatalf("RecordAccess: %v", err)
	}
	if err := RecordAccess(dir, "snap1", "export", "bob", ""); err != nil {
		t.Fatalf("RecordAccess: %v", err)
	}

	log, err := GetAccessLog(dir, "snap1")
	if err != nil {
		t.Fatalf("GetAccessLog: %v", err)
	}
	if len(log.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(log.Events))
	}
	if log.Events[0].Action != "read" || log.Events[0].Actor != "alice" {
		t.Errorf("unexpected first event: %+v", log.Events[0])
	}
	if log.Events[1].Action != "export" || log.Events[1].Note != "" {
		t.Errorf("unexpected second event: %+v", log.Events[1])
	}
}

func TestRecordAccess_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := RecordAccess(dir, "missing", "read", "alice", "")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestGetAccessLog_NoFile(t *testing.T) {
	dir := t.TempDir()
	log, err := GetAccessLog(dir, "snap1")
	if err != nil {
		t.Fatalf("GetAccessLog: %v", err)
	}
	if len(log.Events) != 0 {
		t.Errorf("expected empty log, got %d events", len(log.Events))
	}
}

func TestClearAccessLog(t *testing.T) {
	dir := t.TempDir()
	writeAccessSnapshot(t, dir, "snap1")

	_ = RecordAccess(dir, "snap1", "read", "alice", "")

	if err := ClearAccessLog(dir, "snap1"); err != nil {
		t.Fatalf("ClearAccessLog: %v", err)
	}
	log, _ := GetAccessLog(dir, "snap1")
	if len(log.Events) != 0 {
		t.Errorf("expected empty log after clear, got %d events", len(log.Events))
	}
}

func TestClearAccessLog_NoFile(t *testing.T) {
	dir := t.TempDir()
	if err := ClearAccessLog(dir, "nonexistent"); err != nil {
		t.Errorf("expected no error clearing non-existent log: %v", err)
	}
}
