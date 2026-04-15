package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const pinFileName = "pins.json"

type PinIndex map[string]string // snapshotID -> note

func pinIndexPath(dir string) string {
	return filepath.Join(dir, pinFileName)
}

func LoadPinIndex(dir string) (PinIndex, error) {
	path := pinIndexPath(dir)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return PinIndex{}, nil
	}
	if err != nil {
		return nil, err
	}
	var index PinIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}
	return index, nil
}

func SavePinIndex(dir string, index PinIndex) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pinIndexPath(dir), data, 0644)
}

func PinSnapshot(dir, snapshotID, note string) error {
	if _, err := os.Stat(filepath.Join(dir, snapshotID+".json")); errors.Is(err, os.ErrNotExist) {
		return errors.New("snapshot not found: " + snapshotID)
	}
	index, err := LoadPinIndex(dir)
	if err != nil {
		return err
	}
	index[snapshotID] = note
	return SavePinIndex(dir, index)
}

func UnpinSnapshot(dir, snapshotID string) error {
	index, err := LoadPinIndex(dir)
	if err != nil {
		return err
	}
	if _, ok := index[snapshotID]; !ok {
		return errors.New("snapshot is not pinned: " + snapshotID)
	}
	delete(index, snapshotID)
	return SavePinIndex(dir, index)
}

func IsPinned(dir, snapshotID string) (bool, string, error) {
	index, err := LoadPinIndex(dir)
	if err != nil {
		return false, "", err
	}
	note, ok := index[snapshotID]
	return ok, note, nil
}
