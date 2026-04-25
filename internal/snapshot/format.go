package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FormatConfig defines display formatting preferences for a snapshot.
type FormatConfig struct {
	TimestampLayout string            `json:"timestamp_layout"`
	LevelColors     map[string]string `json:"level_colors"`
	DateFormat      string            `json:"date_format"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

func formatConfigPath(dir, snapshotID string) string {
	return filepath.Join(dir, snapshotID+".format.json")
}

// SetFormatConfig saves a FormatConfig for the given snapshot.
func SetFormatConfig(dir, snapshotID string, cfg FormatConfig) error {
	snaps, err := ListSnapshots(dir)
	if err != nil {
		return fmt.Errorf("list snapshots: %w", err)
	}
	found := false
	for _, s := range snaps {
		if s.ID == snapshotID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("snapshot %q not found", snapshotID)
	}

	now := time.Now().UTC()
	existing, err := GetFormatConfig(dir, snapshotID)
	if err == nil {
		cfg.CreatedAt = existing.CreatedAt
	} else {
		cfg.CreatedAt = now
	}
	cfg.UpdatedAt = now

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal format config: %w", err)
	}
	return os.WriteFile(formatConfigPath(dir, snapshotID), data, 0644)
}

// GetFormatConfig loads the FormatConfig for the given snapshot.
func GetFormatConfig(dir, snapshotID string) (FormatConfig, error) {
	data, err := os.ReadFile(formatConfigPath(dir, snapshotID))
	if err != nil {
		if os.IsNotExist(err) {
			return FormatConfig{}, fmt.Errorf("no format config for snapshot %q", snapshotID)
		}
		return FormatConfig{}, fmt.Errorf("read format config: %w", err)
	}
	var cfg FormatConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return FormatConfig{}, fmt.Errorf("unmarshal format config: %w", err)
	}
	return cfg, nil
}

// RemoveFormatConfig deletes the FormatConfig file for the given snapshot.
func RemoveFormatConfig(dir, snapshotID string) error {
	path := formatConfigPath(dir, snapshotID)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no format config for snapshot %q", snapshotID)
		}
		return fmt.Errorf("remove format config: %w", err)
	}
	return nil
}
