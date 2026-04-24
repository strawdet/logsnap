package snapshot

import (
	"testing"
	"time"
)

func makeSignatureSnapshot(t *testing.T, dir string) *Snapshot {
	t.Helper()
	snap := &Snapshot{
		ID:        "sig-snap-001",
		CreatedAt: time.Now().UTC(),
		Entries: []LogEntry{
			{Timestamp: time.Now().UTC(), Level: "INFO", Message: "service started"},
			{Timestamp: time.Now().UTC(), Level: "WARN", Message: "high memory usage"},
		},
	}
	if err := snap.Save(dir); err != nil {
		t.Fatalf("makeSignatureSnapshot: save: %v", err)
	}
	return snap
}

func TestSignSnapshot_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	snap := makeSignatureSnapshot(t, dir)

	rec, err := SignSnapshot(dir, snap.ID)
	if err != nil {
		t.Fatalf("SignSnapshot: %v", err)
	}
	if rec.Hash == "" {
		t.Error("expected non-empty hash")
	}
	if rec.Algorithm != "sha256" {
		t.Errorf("expected sha256, got %s", rec.Algorithm)
	}
	if rec.SnapshotID != snap.ID {
		t.Errorf("expected snapshot ID %s, got %s", snap.ID, rec.SnapshotID)
	}
}

func TestVerifySnapshot_Valid(t *testing.T) {
	dir := t.TempDir()
	snap := makeSignatureSnapshot(t, dir)

	if _, err := SignSnapshot(dir, snap.ID); err != nil {
		t.Fatalf("SignSnapshot: %v", err)
	}

	ok, rec, err := VerifySnapshot(dir, snap.ID)
	if err != nil {
		t.Fatalf("VerifySnapshot: %v", err)
	}
	if !ok {
		t.Error("expected signature to be valid")
	}
	if rec == nil {
		t.Error("expected non-nil record")
	}
}

func TestVerifySnapshot_Tampered(t *testing.T) {
	dir := t.TempDir()
	snap := makeSignatureSnapshot(t, dir)

	if _, err := SignSnapshot(dir, snap.ID); err != nil {
		t.Fatalf("SignSnapshot: %v", err)
	}

	// Tamper: overwrite the snapshot with different entries.
	snap.Entries = append(snap.Entries, LogEntry{
		Timestamp: time.Now().UTC(),
		Level:     "ERROR",
		Message:   "injected entry",
	})
	if err := snap.Save(dir); err != nil {
		t.Fatalf("save tampered snapshot: %v", err)
	}

	ok, _, err := VerifySnapshot(dir, snap.ID)
	if err != nil {
		t.Fatalf("VerifySnapshot after tamper: %v", err)
	}
	if ok {
		t.Error("expected signature to be invalid after tampering")
	}
}

func TestVerifySnapshot_NoSignature(t *testing.T) {
	dir := t.TempDir()
	snap := makeSignatureSnapshot(t, dir)

	_, _, err := VerifySnapshot(dir, snap.ID)
	if err == nil {
		t.Error("expected error when no signature file exists")
	}
}
