package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeAuditSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	s := &Snapshot{ID: id, Label: "audit-test"}
	data, _ := json.Marshal(s)
	_ = os.WriteFile(filepath.Join(dir, id+".json"), data, 0644)
}

func TestRecordAuditEvent_AndGet(t *testing.T) {
	dir := t.TempDir()
	writeAuditSnapshot(t, dir, "snap1")

	if err := RecordAuditEvent(dir, "snap1", "capture", "initial capture"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := RecordAuditEvent(dir, "snap1", "tag", "tagged as prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	log, err := GetAuditLog(dir, "snap1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(log.Events))
	}
	if log.Events[0].Action != "capture" {
		t.Errorf("expected action 'capture', got %q", log.Events[0].Action)
	}
	if log.Events[1].Detail != "tagged as prod" {
		t.Errorf("expected detail 'tagged as prod', got %q", log.Events[1].Detail)
	}
}

func TestRecordAuditEvent_SnapshotNotFound(t *testing.T) {
	dir := t.TempDir()
	err := RecordAuditEvent(dir, "missing", "capture", "")
	if err == nil {
		t.Fatal("expected error for missing snapshot")
	}
}

func TestGetAuditLog_NoFile(t *testing.T) {
	dir := t.TempDir()
	writeAuditSnapshot(t, dir, "snap2")

	log, err := GetAuditLog(dir, "snap2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Events) != 0 {
		t.Errorf("expected empty log, got %d events", len(log.Events))
	}
}

func TestClearAuditLog(t *testing.T) {
	dir := t.TempDir()
	writeAuditSnapshot(t, dir, "snap3")
	_ = RecordAuditEvent(dir, "snap3", "export", "csv export")

	if err := ClearAuditLog(dir, "snap3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	log, _ := GetAuditLog(dir, "snap3")
	if len(log.Events) != 0 {
		t.Errorf("expected empty log after clear, got %d events", len(log.Events))
	}
}

func TestClearAuditLog_NoFile(t *testing.T) {
	dir := t.TempDir()
	if err := ClearAuditLog(dir, "nonexistent"); err != nil {
		t.Errorf("expected no error clearing absent log, got: %v", err)
	}
}
