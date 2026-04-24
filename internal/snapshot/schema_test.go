package snapshot

import (
	"os"
	"testing"
	"time"
)

func makeSchemaSnapshot(entries []LogEntry) *Snapshot {
	return &Snapshot{
		ID:        "schema-snap-1",
		CreatedAt: time.Now(),
		Entries:   entries,
	}
}

func TestSaveAndLoadSchema(t *testing.T) {
	dir := t.TempDir()
	s := &Schema{
		ID:             "s1",
		Name:           "basic",
		Description:    "basic schema",
		RequiredFields: []string{"message", "level"},
		AllowedLevels:  []string{"info", "error"},
	}
	if err := SaveSchema(dir, s); err != nil {
		t.Fatalf("SaveSchema: %v", err)
	}
	loaded, err := LoadSchema(dir, "basic")
	if err != nil {
		t.Fatalf("LoadSchema: %v", err)
	}
	if loaded.Name != s.Name {
		t.Errorf("expected name %q, got %q", s.Name, loaded.Name)
	}
	if len(loaded.RequiredFields) != 2 {
		t.Errorf("expected 2 required fields, got %d", len(loaded.RequiredFields))
	}
}

func TestLoadSchema_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadSchema(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing schema")
	}
}

func TestDeleteSchema(t *testing.T) {
	dir := t.TempDir()
	s := &Schema{ID: "s2", Name: "temp", RequiredFields: []string{"message"}}
	if err := SaveSchema(dir, s); err != nil {
		t.Fatalf("SaveSchema: %v", err)
	}
	if err := DeleteSchema(dir, "temp"); err != nil {
		t.Fatalf("DeleteSchema: %v", err)
	}
	if err := DeleteSchema(dir, "temp"); err == nil {
		t.Fatal("expected error deleting non-existent schema")
	}
}

func TestValidateSnapshot_NoViolations(t *testing.T) {
	snap := makeSchemaSnapshot([]LogEntry{
		{Message: "started", Level: "info"},
		{Message: "failed", Level: "error"},
	})
	schema := &Schema{
		RequiredFields: []string{"message", "level"},
		AllowedLevels:  []string{"info", "error"},
	}
	violations := ValidateSnapshot(snap, schema)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %+v", len(violations), violations)
	}
}

func TestValidateSnapshot_MissingField(t *testing.T) {
	snap := makeSchemaSnapshot([]LogEntry{
		{Message: "", Level: "info"},
	})
	schema := &Schema{RequiredFields: []string{"message"}}
	violations := ValidateSnapshot(snap, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Field != "message" {
		t.Errorf("expected field 'message', got %q", violations[0].Field)
	}
}

func TestValidateSnapshot_DisallowedLevel(t *testing.T) {
	snap := makeSchemaSnapshot([]LogEntry{
		{Message: "debug msg", Level: "debug"},
	})
	schema := &Schema{
		RequiredFields: []string{"level"},
		AllowedLevels:  []string{"info", "error"},
	}
	violations := ValidateSnapshot(snap, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Field != "level" {
		t.Errorf("expected field 'level', got %q", violations[0].Field)
	}
}

func TestValidateSnapshot_NilSnapshot(t *testing.T) {
	snap := &Snapshot{}
	schema := &Schema{RequiredFields: []string{"message"}}
	violations := ValidateSnapshot(snap, schema)
	if len(violations) != 0 {
		t.Errorf("expected no violations for empty snapshot, got %d", len(violations))
	}
}

func init() { _ = os.Getenv("") }
