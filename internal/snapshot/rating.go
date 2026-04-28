package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Rating struct {
	Score     int       `json:"score"`      // 1-5
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ratingPath(dir, id string) string {
	return filepath.Join(dir, id+".rating.json")
}

func SetRating(dir, id string, score int, comment string) error {
	if score < 1 || score > 5 {
		return fmt.Errorf("score must be between 1 and 5, got %d", score)
	}

	snapshotFile := filepath.Join(dir, id+".json")
	if _, err := os.Stat(snapshotFile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("snapshot %q not found", id)
	}

	now := time.Now().UTC()
	r := Rating{
		Score:     score,
		Comment:  comment,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Preserve original CreatedAt if rating already exists
	existing, err := GetRating(dir, id)
	if err == nil {
		r.CreatedAt = existing.CreatedAt
	}

	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal rating: %w", err)
	}
	return os.WriteFile(ratingPath(dir, id), data, 0644)
}

func GetRating(dir, id string) (*Rating, error) {
	path := ratingPath(dir, id)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("no rating for snapshot %q", id)
	}
	if err != nil {
		return nil, fmt.Errorf("read rating: %w", err)
	}
	var r Rating
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("unmarshal rating: %w", err)
	}
	return &r, nil
}

func RemoveRating(dir, id string) error {
	path := ratingPath(dir, id)
	if err := os.Remove(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("no rating found for snapshot %q", id)
	} else if err != nil {
		return fmt.Errorf("remove rating: %w", err)
	}
	return nil
}
