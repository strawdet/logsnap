package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type AliasIndex map[string]string // alias -> snapshotID

func aliasIndexPath(dir string) string {
	return filepath.Join(dir, "aliases.json")
}

func LoadAliasIndex(dir string) (AliasIndex, error) {
	path := aliasIndexPath(dir)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return AliasIndex{}, nil
	}
	if err != nil {
		return nil, err
	}
	var idx AliasIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return idx, nil
}

func SaveAliasIndex(dir string, idx AliasIndex) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(aliasIndexPath(dir), data, 0644)
}

func SetAlias(dir, alias, snapshotID string) error {
	path := filepath.Join(dir, snapshotID+".json")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return errors.New("snapshot not found: " + snapshotID)
	}
	idx, err := LoadAliasIndex(dir)
	if err != nil {
		return err
	}
	idx[alias] = snapshotID
	return SaveAliasIndex(dir, idx)
}

func ResolveAlias(dir, alias string) (string, error) {
	idx, err := LoadAliasIndex(dir)
	if err != nil {
		return "", err
	}
	id, ok := idx[alias]
	if !ok {
		return "", errors.New("alias not found: " + alias)
	}
	return id, nil
}

func RemoveAlias(dir, alias string) error {
	idx, err := LoadAliasIndex(dir)
	if err != nil {
		return err
	}
	if _, ok := idx[alias]; !ok {
		return errors.New("alias not found: " + alias)
	}
	delete(idx, alias)
	return SaveAliasIndex(dir, idx)
}

func ListAliases(dir string) (AliasIndex, error) {
	return LoadAliasIndex(dir)
}

// AliasesForSnapshot returns all aliases that point to the given snapshotID.
func AliasesForSnapshot(dir, snapshotID string) ([]string, error) {
	idx, err := LoadAliasIndex(dir)
	if err != nil {
		return nil, err
	}
	var aliases []string
	for alias, id := range idx {
		if id == snapshotID {
			aliases = append(aliases, alias)
		}
	}
	return aliases, nil
}
