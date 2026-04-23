package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Group represents a named collection of snapshot IDs.
type Group struct {
	Name      string    `json:"name"`
	Snapshots []string  `json:"snapshots"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func groupIndexPath(dir string) string {
	return filepath.Join(dir, "group_index.json")
}

// LoadGroupIndex reads all groups from disk.
func LoadGroupIndex(dir string) (map[string]*Group, error) {
	path := groupIndexPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]*Group{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read group index: %w", err)
	}
	var index map[string]*Group
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("parse group index: %w", err)
	}
	return index, nil
}

// SaveGroupIndex writes the group index to disk.
func SaveGroupIndex(dir string, index map[string]*Group) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal group index: %w", err)
	}
	return os.WriteFile(groupIndexPath(dir), data, 0644)
}

// AddToGroup adds a snapshot ID to a named group, creating it if needed.
func AddToGroup(dir, groupName, snapshotID string) error {
	snapshotFile := filepath.Join(dir, snapshotID+".json")
	if _, err := os.Stat(snapshotFile); os.IsNotExist(err) {
		return fmt.Errorf("snapshot %q not found", snapshotID)
	}
	index, err := LoadGroupIndex(dir)
	if err != nil {
		return err
	}
	g, ok := index[groupName]
	if !ok {
		g = &Group{Name: groupName, CreatedAt: time.Now()}
	}
	for _, id := range g.Snapshots {
		if id == snapshotID {
			return nil // already present
		}
	}
	g.Snapshots = append(g.Snapshots, snapshotID)
	g.UpdatedAt = time.Now()
	index[groupName] = g
	return SaveGroupIndex(dir, index)
}

// RemoveFromGroup removes a snapshot ID from a named group.
func RemoveFromGroup(dir, groupName, snapshotID string) error {
	index, err := LoadGroupIndex(dir)
	if err != nil {
		return err
	}
	g, ok := index[groupName]
	if !ok {
		return fmt.Errorf("group %q not found", groupName)
	}
	updated := g.Snapshots[:0]
	for _, id := range g.Snapshots {
		if id != snapshotID {
			updated = append(updated, id)
		}
	}
	g.Snapshots = updated
	g.UpdatedAt = time.Now()
	index[groupName] = g
	return SaveGroupIndex(dir, index)
}

// GetGroup returns the group with the given name.
func GetGroup(dir, groupName string) (*Group, error) {
	index, err := LoadGroupIndex(dir)
	if err != nil {
		return nil, err
	}
	g, ok := index[groupName]
	if !ok {
		return nil, fmt.Errorf("group %q not found", groupName)
	}
	return g, nil
}

// ListGroups returns all group names.
func ListGroups(dir string) ([]string, error) {
	index, err := LoadGroupIndex(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(index))
	for name := range index {
		names = append(names, name)
	}
	return names, nil
}
