package snapshot

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// SignatureIndexEntry holds metadata about a signed snapshot.
type SignatureIndexEntry struct {
	SnapshotID string    `json:"snapshot_id"`
	SignedAt   time.Time `json:"signed_at"`
	Signer     string    `json:"signer"`
	Valid      bool      `json:"valid"`
}

// SignatureIndex maps snapshot IDs to their signature metadata.
type SignatureIndex struct {
	Entries map[string]SignatureIndexEntry `json:"entries"`
}

func signatureIndexPath(dir string) string {
	return filepath.Join(dir, ".signature_index.json")
}

// LoadSignatureIndex loads the signature index from the given directory.
// Returns an empty index if the file does not exist.
func LoadSignatureIndex(dir string) (*SignatureIndex, error) {
	path := signatureIndexPath(dir)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &SignatureIndex{Entries: make(map[string]SignatureIndexEntry)}, nil
		}
		return nil, err
	}
	var idx SignatureIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	if idx.Entries == nil {
		idx.Entries = make(map[string]SignatureIndexEntry)
	}
	return &idx, nil
}

// SaveSignatureIndex persists the signature index to the given directory.
func SaveSignatureIndex(dir string, idx *SignatureIndex) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(signatureIndexPath(dir), data, 0644)
}

// RegisterSignature adds or updates a signature entry in the index.
func RegisterSignature(dir, snapshotID, signer string, valid bool) error {
	idx, err := LoadSignatureIndex(dir)
	if err != nil {
		return err
	}
	idx.Entries[snapshotID] = SignatureIndexEntry{
		SnapshotID: snapshotID,
		SignedAt:   time.Now().UTC(),
		Signer:     signer,
		Valid:      valid,
	}
	return SaveSignatureIndex(dir, idx)
}

// DeregisterSignature removes a snapshot's signature entry from the index.
func DeregisterSignature(dir, snapshotID string) error {
	idx, err := LoadSignatureIndex(dir)
	if err != nil {
		return err
	}
	delete(idx.Entries, snapshotID)
	return SaveSignatureIndex(dir, idx)
}

// ListSignedSnapshots returns all signature index entries.
func ListSignedSnapshots(dir string) ([]SignatureIndexEntry, error) {
	idx, err := LoadSignatureIndex(dir)
	if err != nil {
		return nil, err
	}
	result := make([]SignatureIndexEntry, 0, len(idx.Entries))
	for _, entry := range idx.Entries {
		result = append(result, entry)
	}
	return result, nil
}
