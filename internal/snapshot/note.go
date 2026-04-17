package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Note struct {
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func notePath(dir, id string) string {
	return filepath.Join(dir, id+".note.json")
}

func AddNote(dir, id, text string) error {
	snapshotFile := filepath.Join(dir, id+".json")
	if _, err := os.Stat(snapshotFile); os.IsNotExist(err) {
		return fmt.Errorf("snapshot %s not found", id)
	}
	now := time.Now()
	note := Note{Text: text, CreatedAt: now, UpdatedAt: now}
	existing, err := GetNote(dir, id)
	if err == nil {
		note.CreatedAt = existing.CreatedAt
	}
	data, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(notePath(dir, id), data, 0644)
}

func GetNote(dir, id string) (*Note, error) {
	data, err := os.ReadFile(notePath(dir, id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no note for snapshot %s", id)
		}
		return nil, err
	}
	var note Note
	if err := json.Unmarshal(data, &note); err != nil {
		return nil, err
	}
	return &note, nil
}

func RemoveNote(dir, id string) error {
	p := notePath(dir, id)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("no note for snapshot %s", id)
	}
	return os.Remove(p)
}
