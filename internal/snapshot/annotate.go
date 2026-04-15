package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Annotation holds a user-defined note attached to a snapshot.
type Annotation struct {
	SnapshotID string    `json:"snapshot_id"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
}

func annotationPath(dir, snapshotID string) string {
	return filepath.Join(dir, snapshotID+".note.json")
}

// AddAnnotation writes a note for the given snapshot ID.
func AddAnnotation(dir, snapshotID, note string) error {
	path := filepath.Join(dir, snapshotID+".json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("snapshot %q not found", snapshotID)
	}

	a := Annotation{
		SnapshotID: snapshotID,
		Note:       note,
		CreatedAt:  time.Now().UTC(),
	}

	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal annotation: %w", err)
	}

	return os.WriteFile(annotationPath(dir, snapshotID), data, 0644)
}

// GetAnnotation reads the note for the given snapshot ID.
// Returns nil, nil if no annotation exists.
func GetAnnotation(dir, snapshotID string) (*Annotation, error) {
	p := annotationPath(dir, snapshotID)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read annotation: %w", err)
	}

	var a Annotation
	if err := json.Unmarshal(data, &a); err != nil {
		return nil, fmt.Errorf("unmarshal annotation: %w", err)
	}
	return &a, nil
}

// RemoveAnnotation deletes the annotation file for a snapshot.
func RemoveAnnotation(dir, snapshotID string) error {
	p := annotationPath(dir, snapshotID)
	if err := os.Remove(p); os.IsNotExist(err) {
		return fmt.Errorf("annotation for %q not found", snapshotID)
	} else if err != nil {
		return fmt.Errorf("remove annotation: %w", err)
	}
	return nil
}
