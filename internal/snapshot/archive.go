package snapshot

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ArchiveSnapshots zips the given snapshot IDs into a single archive file.
func ArchiveSnapshots(dir string, ids []string, outPath string) error {
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create archive: %w", err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	for _, id := range ids {
		snap, err := Load(dir, id)
		if err != nil {
			return fmt.Errorf("load snapshot %s: %w", id, err)
		}
		data, err := json.Marshal(snap)
		if err != nil {
			return fmt.Errorf("marshal snapshot %s: %w", id, err)
		}
		w, err := zw.Create(id + ".json")
		if err != nil {
			return fmt.Errorf("zip entry %s: %w", id, err)
		}
		if _, err := w.Write(data); err != nil {
			return fmt.Errorf("write entry %s: %w", id, err)
		}
	}
	return nil
}

// UnarchiveSnapshots extracts snapshots from a zip archive into dir.
func UnarchiveSnapshots(dir string, archivePath string) ([]string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, fmt.Errorf("open archive: %w", err)
	}
	defer r.Close()

	var ids []string
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
			return nil, fmt.Errorf("unmarshal entry %s: %w", f.Name, err)
		}
		out := filepath.Join(dir, snap.ID+".json")
		if err := os.WriteFile(out, data, 0644); err != nil {
			return nil, fmt.Errorf("write snapshot %s: %w", snap.ID, err)
		}
		ids = append(ids, snap.ID)
	}
	return ids, nil
}
