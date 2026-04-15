package snapshot

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a captured log snapshot at a point in time.
type Snapshot struct {
	ID        string            `json:"id"`
	Label     string            `json:"label"`
	CreatedAt time.Time         `json:"created_at"`
	Entries   []LogEntry        `json:"entries"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// LogEntry represents a single structured log line.
type LogEntry struct {
	Timestamp string            `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
}

// New creates a new Snapshot with a generated ID.
func New(label string, entries []LogEntry, meta map[string]string) *Snapshot {
	s := &Snapshot{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Entries:   entries,
		Meta:      meta,
	}
	s.ID = generateID(s)
	return s
}

// Save writes the snapshot as JSON to the given file path.
func (s *Snapshot) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads a snapshot from a JSON file.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read snapshot file: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal snapshot: %w", err)
	}
	return &s, nil
}

func generateID(s *Snapshot) string {
	h := sha256.New()
	h.Write([]byte(s.Label))
	h.Write([]byte(s.CreatedAt.String()))
	return fmt.Sprintf("%x", h.Sum(nil))[:12]
}
