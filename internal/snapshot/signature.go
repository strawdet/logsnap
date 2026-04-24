package snapshot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SignatureRecord holds the computed signature for a snapshot.
type SignatureRecord struct {
	SnapshotID string    `json:"snapshot_id"`
	Hash       string    `json:"hash"`
	Algorithm  string    `json:"algorithm"`
	SignedAt   time.Time `json:"signed_at"`
}

func signaturePath(dir, snapshotID string) string {
	return filepath.Join(dir, snapshotID+".sig.json")
}

// SignSnapshot computes a SHA-256 hash over the snapshot's entries and saves it.
func SignSnapshot(dir, snapshotID string) (*SignatureRecord, error) {
	snap, err := Load(dir, snapshotID)
	if err != nil {
		return nil, fmt.Errorf("sign: load snapshot: %w", err)
	}

	data, err := json.Marshal(snap.Entries)
	if err != nil {
		return nil, fmt.Errorf("sign: marshal entries: %w", err)
	}

	sum := sha256.Sum256(data)
	rec := &SignatureRecord{
		SnapshotID: snapshotID,
		Hash:       hex.EncodeToString(sum[:]),
		Algorithm:  "sha256",
		SignedAt:   time.Now().UTC(),
	}

	out, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("sign: marshal record: %w", err)
	}
	if err := os.WriteFile(signaturePath(dir, snapshotID), out, 0644); err != nil {
		return nil, fmt.Errorf("sign: write file: %w", err)
	}
	return rec, nil
}

// VerifySnapshot loads the stored signature and re-computes the hash to confirm integrity.
func VerifySnapshot(dir, snapshotID string) (bool, *SignatureRecord, error) {
	recBytes, err := os.ReadFile(signaturePath(dir, snapshotID))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, fmt.Errorf("verify: no signature found for %s", snapshotID)
		}
		return false, nil, fmt.Errorf("verify: read sig: %w", err)
	}

	var rec SignatureRecord
	if err := json.Unmarshal(recBytes, &rec); err != nil {
		return false, nil, fmt.Errorf("verify: parse sig: %w", err)
	}

	snap, err := Load(dir, snapshotID)
	if err != nil {
		return false, &rec, fmt.Errorf("verify: load snapshot: %w", err)
	}

	data, err := json.Marshal(snap.Entries)
	if err != nil {
		return false, &rec, fmt.Errorf("verify: marshal entries: %w", err)
	}

	sum := sha256.Sum256(data)
	actual := hex.EncodeToString(sum[:])
	return actual == rec.Hash, &rec, nil
}
