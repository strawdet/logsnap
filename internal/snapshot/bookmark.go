package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type BookmarkIndex map[string]string // name -> snapshotID

func bookmarkPath(dir string) string {
	return filepath.Join(dir, "bookmarks.json")
}

func LoadBookmarkIndex(dir string) (BookmarkIndex, error) {
	p := bookmarkPath(dir)
	data, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return make(BookmarkIndex), nil
	}
	if err != nil {
		return nil, err
	}
	var idx BookmarkIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return idx, nil
}

func SaveBookmarkIndex(dir string, idx BookmarkIndex) error {
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(bookmarkPath(dir), data, 0644)
}

func AddBookmark(dir, name, snapshotID string) error {
	if _, err := Load(dir, snapshotID); err != nil {
		return errors.New("snapshot not found: " + snapshotID)
	}
	idx, err := LoadBookmarkIndex(dir)
	if err != nil {
		return err
	}
	idx[name] = snapshotID
	return SaveBookmarkIndex(dir, idx)
}

func ResolveBookmark(dir, name string) (string, error) {
	idx, err := LoadBookmarkIndex(dir)
	if err != nil {
		return "", err
	}
	id, ok := idx[name]
	if !ok {
		return "", errors.New("bookmark not found: " + name)
	}
	return id, nil
}

func RemoveBookmark(dir, name string) error {
	idx, err := LoadBookmarkIndex(dir)
	if err != nil {
		return err
	}
	if _, ok := idx[name]; !ok {
		return errors.New("bookmark not found: " + name)
	}
	delete(idx, name)
	return SaveBookmarkIndex(dir, idx)
}

func ListBookmarks(dir string) (BookmarkIndex, error) {
	return LoadBookmarkIndex(dir)
}
