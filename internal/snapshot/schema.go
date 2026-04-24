package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Schema defines a validation schema for snapshot log entries.
type Schema struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	RequiredFields []string       `json:"required_fields"`
	AllowedLevels  []string       `json:"allowed_levels,omitempty"`
}

// SchemaViolation describes a single validation failure.
type SchemaViolation struct {
	EntryIndex int    `json:"entry_index"`
	Field      string `json:"field"`
	Reason     string `json:"reason"`
}

func schemaPath(dir, name string) string {
	return filepath.Join(dir, "schemas", name+".json")
}

// SaveSchema persists a schema to disk.
func SaveSchema(dir string, schema *Schema) error {
	p := schemaPath(dir, schema.Name)
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return fmt.Errorf("create schema dir: %w", err)
	}
	schema.CreatedAt = time.Now()
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal schema: %w", err)
	}
	return os.WriteFile(p, data, 0644)
}

// LoadSchema reads a schema from disk by name.
func LoadSchema(dir, name string) (*Schema, error) {
	data, err := os.ReadFile(schemaPath(dir, name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("schema %q not found", name)
		}
		return nil, fmt.Errorf("read schema: %w", err)
	}
	var s Schema
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal schema: %w", err)
	}
	return &s, nil
}

// DeleteSchema removes a schema file from disk.
func DeleteSchema(dir, name string) error {
	p := schemaPath(dir, name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("schema %q not found", name)
	}
	return os.Remove(p)
}

// ValidateSnapshot checks a snapshot's entries against the given schema.
func ValidateSnapshot(snap *Snapshot, schema *Schema) []SchemaViolation {
	var violations []SchemaViolation
	allowedSet := make(map[string]bool, len(schema.AllowedLevels))
	for _, l := range schema.AllowedLevels {
		allowedSet[l] = true
	}
	for i, entry := range snap.Entries {
		for _, field := range schema.RequiredFields {
			switch field {
			case "message":
				if entry.Message == "" {
					violations = append(violations, SchemaViolation{EntryIndex: i, Field: field, Reason: "missing required field"})
				}
			case "level":
				if entry.Level == "" {
					violations = append(violations, SchemaViolation{EntryIndex: i, Field: field, Reason: "missing required field"})
				}
			}
		}
		if len(allowedSet) > 0 && entry.Level != "" && !allowedSet[entry.Level] {
			violations = append(violations, SchemaViolation{EntryIndex: i, Field: "level", Reason: fmt.Sprintf("level %q not in allowed set", entry.Level)})
		}
	}
	return violations
}
