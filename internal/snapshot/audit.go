package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuditEvent represents a single recorded action on a snapshot.
type AuditEvent struct {
	Timestamp  time.Time `json:"timestamp"`
	Action     string    `json:"action"`
	SnapshotID string    `json:"snapshot_id"`
	Detail     string    `json:"detail,omitempty"`
}

// AuditLog holds all events for a snapshot.
type AuditLog struct {
	Events []AuditEvent `json:"events"`
}

func auditPath(dir, snapshotID string) string {
	return filepath.Join(dir, snapshotID+".audit.json")
}

// RecordAuditEvent appends an event to the audit log for a snapshot.
func RecordAuditEvent(dir, snapshotID, action, detail string) error {
	path := filepath.Join(dir, snapshotID+".json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("snapshot %s not found", snapshotID)
	}

	log, _ := GetAuditLog(dir, snapshotID)
	log.Events = append(log.Events, AuditEvent{
		Timestamp:  time.Now().UTC(),
		Action:     action,
		SnapshotID: snapshotID,
		Detail:     detail,
	})

	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal audit log: %w", err)
	}
	return os.WriteFile(auditPath(dir, snapshotID), data, 0644)
}

// GetAuditLog loads the audit log for a snapshot. Returns an empty log if none exists.
func GetAuditLog(dir, snapshotID string) (*AuditLog, error) {
	p := auditPath(dir, snapshotID)
	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return &AuditLog{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read audit log: %w", err)
	}
	var log AuditLog
	if err := json.Unmarshal(data, &log); err != nil {
		return nil, fmt.Errorf("parse audit log: %w", err)
	}
	return &log, nil
}

// ClearAuditLog removes the audit log file for a snapshot.
func ClearAuditLog(dir, snapshotID string) error {
	p := auditPath(dir, snapshotID)
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove audit log: %w", err)
	}
	return nil
}
