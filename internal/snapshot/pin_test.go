package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func writePinSnapshot(t *testing.T, dir, id string) {
	t.Helper()
	f := filepath.Join(dir, id+".json")
	if err := os.WriteFile(f, []byte(`{"id":"`+id+`"}`), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestPinSnapshot_AndIsPinned(t *testing.T) {
	dir := t.TempDir()
	writePinSnapshot(t, dir, "snap-001")

	if err := PinSnapshot(dir, "snap-001", "important release"); err != nil {
		t.Fatalf("PinSnapshot failed: %v", err)
	}

	ok, note, err := IsPinned(dir, "snap-001")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("expected snapshot to be pinned")
	}
	if note != "important release" {
		t.Errorf("expected note 'important release', got %q", note)
	}
}

func TestPinSnapshot_NotFound(t *testing.T) {
	dir := t.TempDir()
	err := PinSnapshot(dir, "nonexistent", "note")
	if err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestUnpinSnapshot(t *testing.T) {
	dir := t.TempDir()
	writePinSnapshot(t, dir, "snap-002")

	_ = PinSnapshot(dir, "snap-002", "keep")

	if err := UnpinSnapshot(dir, "snap-002"); err != nil {
		t.Fatalf("UnpinSnapshot failed: %v", err)
	}

	ok, _, err := IsPinned(dir, "snap-002")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("expected snapshot to be unpinned")
	}
}

func TestUnpinSnapshot_NotPinned(t *testing.T) {
	dir := t.TempDir()
	writePinSnapshot(t, dir, "snap-003")

	err := UnpinSnapshot(dir, "snap-003")
	if err == nil {
		t.Error("expected error when unpinning a non-pinned snapshot")
	}
}

func TestLoadPinIndex_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	index, err := LoadPinIndex(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(index) != 0 {
		t.Errorf("expected empty index, got %d entries", len(index))
	}
}
