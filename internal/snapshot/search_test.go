package snapshot

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeSearchSnapshot(t *testing.T, dir, id string, tags []string, labels map[string]string, createdAt time.Time) *Snapshot {
	t.Helper()
	snap := &Snapshot{
		ID:        id,
		CreatedAt: createdAt,
		Tags:      tags,
		Labels:    labels,
		Entries:   sampleEntries(),
	}
	path := filepath.Join(dir, id+".json")
	err := snap.Save(path)
	require.NoError(t, err)
	return snap
}

func TestSearch_ByTag(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	makeSearchSnapshot(t, dir, "snap-a", []string{"production"}, nil, now)
	makeSearchSnapshot(t, dir, "snap-b", []string{"staging"}, nil, now)

	results, err := Search(dir, SearchFilter{Tag: "production"})
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "snap-a", results[0].Snapshot.ID)
}

func TestSearch_BySince(t *testing.T) {
	dir := t.TempDir()
	old := time.Now().Add(-48 * time.Hour)
	recent := time.Now().Add(-1 * time.Hour)
	makeSearchSnapshot(t, dir, "snap-old", nil, nil, old)
	makeSearchSnapshot(t, dir, "snap-new", nil, nil, recent)

	cutoff := time.Now().Add(-24 * time.Hour)
	results, err := Search(dir, SearchFilter{Since: &cutoff})
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "snap-new", results[0].Snapshot.ID)
}

func TestSearch_ByLabel(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	makeSearchSnapshot(t, dir, "snap-x", nil, map[string]string{"env": "prod"}, now)
	makeSearchSnapshot(t, dir, "snap-y", nil, map[string]string{"env": "dev"}, now)

	results, err := Search(dir, SearchFilter{LabelKey: "env", LabelVal: "prod"})
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "snap-x", results[0].Snapshot.ID)
}

func TestSearch_EmptyFilter(t *testing.T) {
	dir := t.TempDir()
	now := time.Now()
	makeSearchSnapshot(t, dir, "snap-1", nil, nil, now)
	makeSearchSnapshot(t, dir, "snap-2", nil, nil, now)

	results, err := Search(dir, SearchFilter{})
	require.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestSearch_NonExistentDir(t *testing.T) {
	_, err := Search(filepath.Join(os.TempDir(), "no-such-dir-logsnap"), SearchFilter{})
	assert.NoError(t, err) // ListSnapshots returns empty for missing dir
}
