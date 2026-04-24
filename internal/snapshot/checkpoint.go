package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Checkpoint represents a named point-in-time marker referencing a snapshot.
type Checkpoint struct {
	Name        string    `json:"name"`
	SnapshotID  string    `json:"snapshot_id"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CheckpointIndex maps checkpoint names to Checkpoint entries.
type CheckpointIndex map[string]Checkpoint

func checkpointPath(dir string) string {
	return filepath.Join(dir, ".checkpoints.json")
}

// LoadCheckpointIndex loads the checkpoint index from dir.
func LoadCheckpointIndex(dir string) (CheckpointIndex, error) {
	path := checkpointPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(CheckpointIndex), nil
	}
	if err != nil {
		return nil, err
	}
	var idx CheckpointIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return idx, nil
}

// SaveCheckpointIndex persists the checkpoint index to dir.
func SaveCheckpointIndex(dir string, idx CheckpointIndex) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(checkpointPath(dir), data, 0644)
}

// SetCheckpoint creates or updates a named checkpoint pointing to snapshotID.
func SetCheckpoint(dir, name, snapshotID, description string) error {
	snapshotFile := filepath.Join(dir, snapshotID+".json")
	if _, err := os.Stat(snapshotFile); os.IsNotExist(err) {
		return fmt.Errorf("snapshot %q not found", snapshotID)
	}
	idx, err := LoadCheckpointIndex(dir)
	if err != nil {
		return err
	}
	idx[name] = Checkpoint{
		Name:        name,
		SnapshotID:  snapshotID,
		Description: description,
		CreatedAt:   time.Now().UTC(),
	}
	return SaveCheckpointIndex(dir, idx)
}

// ResolveCheckpoint returns the snapshot ID associated with name.
func ResolveCheckpoint(dir, name string) (string, error) {
	idx, err := LoadCheckpointIndex(dir)
	if err != nil {
		return "", err
	}
	cp, ok := idx[name]
	if !ok {
		return "", fmt.Errorf("checkpoint %q not found", name)
	}
	return cp.SnapshotID, nil
}

// RemoveCheckpoint deletes a checkpoint by name.
func RemoveCheckpoint(dir, name string) error {
	idx, err := LoadCheckpointIndex(dir)
	if err != nil {
		return err
	}
	if _, ok := idx[name]; !ok {
		return fmt.Errorf("checkpoint %q not found", name)
	}
	delete(idx, name)
	return SaveCheckpointIndex(dir, idx)
}

// ListCheckpoints returns all checkpoints in the index.
func ListCheckpoints(dir string) ([]Checkpoint, error) {
	idx, err := LoadCheckpointIndex(dir)
	if err != nil {
		return nil, err
	}
	result := make([]Checkpoint, 0, len(idx))
	for _, cp := range idx {
		result = append(result, cp)
	}
	return result, nil
}
