package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type AccessEvent struct {
	SnapshotID string    `json:"snapshot_id"`
	Action     string    `json:"action"`
	Actor      string    `json:"actor"`
	At         time.Time `json:"at"`
	Note       string    `json:"note,omitempty"`
}

type AccessLog struct {
	Events []AccessEvent `json:"events"`
}

func accessPath(dir, snapshotID string) string {
	return filepath.Join(dir, snapshotID+".access.json")
}

func RecordAccess(dir, snapshotID, action, actor, note string) error {
	snapshotFile := filepath.Join(dir, snapshotID+".json")
	if _, err := os.Stat(snapshotFile); os.IsNotExist(err) {
		return fmt.Errorf("snapshot %s not found", snapshotID)
	}

	log, _ := GetAccessLog(dir, snapshotID)
	log.Events = append(log.Events, AccessEvent{
		SnapshotID: snapshotID,
		Action:     action,
		Actor:      actor,
		At:         time.Now().UTC(),
		Note:       note,
	})

	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal access log: %w", err)
	}
	return os.WriteFile(accessPath(dir, snapshotID), data, 0644)
}

func GetAccessLog(dir, snapshotID string) (*AccessLog, error) {
	path := accessPath(dir, snapshotID)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &AccessLog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read access log: %w", err)
	}
	var log AccessLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("unmarshal access log: %w", err)
	}
	return &log, nil
}

func ClearAccessLog(dir, snapshotID string) error {
	path := accessPath(dir, snapshotID)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("clear access log: %w", err)
	}
	return nil
}
