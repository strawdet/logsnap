package snapshot

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
)

// ArchiveInfo holds metadata about a snapshot inside an archive.
type ArchiveInfo struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// ListArchive returns metadata for all snapshots in a zip archive without extracting them.
func ListArchive(archivePath string) ([]ArchiveInfo, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, fmt.Errorf("open archive: %w", err)
	}
	defer r.Close()

	var infos []ArchiveInfo
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("open entry %s: %w", f.Name, err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("read entry %s: %w", f.Name, err)
		}
		var snap Snapshot
		if err := json.Unmarshal(data, &snap); err != nil {
			return nil, fmt.Errorf("parse entry %s: %w", f.Name, err)
		}
		infos = append(infos, ArchiveInfo{ID: snap.ID, Label: snap.Label})
	}
	return infos, nil
}
