package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// TriggerCondition defines when an automatic snapshot should be taken.
type TriggerCondition struct {
	ErrorRateThreshold float64 `json:"error_rate_threshold,omitempty"` // 0.0–1.0
	MinEntries         int     `json:"min_entries,omitempty"`
	Label              string  `json:"label,omitempty"`
}

// Trigger represents a saved auto-capture trigger rule.
type Trigger struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Condition TriggerCondition `json:"condition"`
	CreatedAt time.Time        `json:"created_at"`
}

func triggerPath(dir, name string) string {
	return filepath.Join(dir, "triggers", name+".json")
}

// SaveTrigger persists a named trigger rule to disk.
func SaveTrigger(dir, name string, cond TriggerCondition) (*Trigger, error) {
	if name == "" {
		return nil, fmt.Errorf("trigger name must not be empty")
	}
	t := &Trigger{
		ID:        generateID(),
		Name:      name,
		Condition: cond,
		CreatedAt: time.Now().UTC(),
	}
	path := triggerPath(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create trigger dir: %w", err)
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal trigger: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return nil, fmt.Errorf("write trigger: %w", err)
	}
	return t, nil
}

// LoadTrigger reads a named trigger from disk.
func LoadTrigger(dir, name string) (*Trigger, error) {
	path := triggerPath(dir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("trigger %q not found", name)
		}
		return nil, fmt.Errorf("read trigger: %w", err)
	}
	var t Trigger
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("unmarshal trigger: %w", err)
	}
	return &t, nil
}

// DeleteTrigger removes a named trigger from disk.
func DeleteTrigger(dir, name string) error {
	path := triggerPath(dir, name)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("trigger %q not found", name)
		}
		return fmt.Errorf("delete trigger: %w", err)
	}
	return nil
}

// EvaluateTrigger returns true when the snapshot satisfies the trigger condition.
func EvaluateTrigger(snap *Snapshot, cond TriggerCondition) bool {
	if snap == nil || len(snap.Entries) == 0 {
		return false
	}
	if cond.MinEntries > 0 && len(snap.Entries) < cond.MinEntries {
		return false
	}
	if cond.ErrorRateThreshold > 0 {
		errCount := 0
		for _, e := range snap.Entries {
			if e.Level == "error" || e.Level == "ERROR" {
				errCount++
			}
		}
		rate := float64(errCount) / float64(len(snap.Entries))
		if rate < cond.ErrorRateThreshold {
			return false
		}
	}
	return true
}
