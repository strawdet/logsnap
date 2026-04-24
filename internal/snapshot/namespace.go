package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type NamespaceIndex struct {
	Namespaces map[string][]string `json:"namespaces"` // namespace -> snapshot IDs
	UpdatedAt  time.Time           `json:"updated_at"`
}

func namespacePath(dir string) string {
	return filepath.Join(dir, "namespace_index.json")
}

func LoadNamespaceIndex(dir string) (*NamespaceIndex, error) {
	path := namespacePath(dir)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &NamespaceIndex{Namespaces: make(map[string][]string)}, nil
	}
	if err != nil {
		return nil, err
	}
	var idx NamespaceIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	if idx.Namespaces == nil {
		idx.Namespaces = make(map[string][]string)
	}
	return &idx, nil
}

func SaveNamespaceIndex(dir string, idx *NamespaceIndex) error {
	idx.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(namespacePath(dir), data, 0644)
}

func AddToNamespace(dir, namespace, snapshotID string) error {
	if _, err := Load(dir, snapshotID); err != nil {
		return errors.New("snapshot not found: " + snapshotID)
	}
	idx, err := LoadNamespaceIndex(dir)
	if err != nil {
		return err
	}
	for _, id := range idx.Namespaces[namespace] {
		if id == snapshotID {
			return nil
		}
	}
	idx.Namespaces[namespace] = append(idx.Namespaces[namespace], snapshotID)
	return SaveNamespaceIndex(dir, idx)
}

func RemoveFromNamespace(dir, namespace, snapshotID string) error {
	idx, err := LoadNamespaceIndex(dir)
	if err != nil {
		return err
	}
	ids := idx.Namespaces[namespace]
	filtered := ids[:0]
	for _, id := range ids {
		if id != snapshotID {
			filtered = append(filtered, id)
		}
	}
	idx.Namespaces[namespace] = filtered
	return SaveNamespaceIndex(dir, idx)
}

func ListNamespaces(dir string) ([]string, error) {
	idx, err := LoadNamespaceIndex(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(idx.Namespaces))
	for ns := range idx.Namespaces {
		names = append(names, ns)
	}
	sort.Strings(names)
	return names, nil
}

func GetNamespaceSnapshots(dir, namespace string) ([]string, error) {
	idx, err := LoadNamespaceIndex(dir)
	if err != nil {
		return nil, err
	}
	return idx.Namespaces[namespace], nil
}
