package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logsnap/internal/snapshot"
)

func TestSaveAndLoadCompareResult(t *testing.T) {
	dir := t.TempDir()

	result, err := snapshot.SaveCompareResult(dir, "deploy-v1-v2", "abc123", "def456")
	if err != nil {
		t.Fatalf("SaveCompareResult: %v", err)
	}

	if result.Name != "deploy-v1-v2" {
		t.Errorf("expected name deploy-v1-v2, got %s", result.Name)
	}

	loaded, err := snapshot.LoadCompareResult(dir, "deploy-v1-v2")
	if err != nil {
		t.Fatalf("LoadCompareResult: %v", err)
	}

	if loaded.BaseID != "abc123" {
		t.Errorf("expected base abc123, got %s", loaded.BaseID)
	}
	if loaded.TargetID != "def456" {
		t.Errorf("expected target def456, got %s", loaded.TargetID)
	}
}

func TestLoadCompareResult_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := snapshot.LoadCompareResult(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing compare result")
	}
}

func TestListCompareResults(t *testing.T) {
	dir := t.TempDir()

	for _, name := range []string{"alpha", "beta", "gamma"} {
		_, err := snapshot.SaveCompareResult(dir, name, "id1", "id2")
		if err != nil {
			t.Fatalf("SaveCompareResult %s: %v", name, err)
		}
	}

	names, err := snapshot.ListCompareResults(dir)
	if err != nil {
		t.Fatalf("ListCompareResults: %v", err)
	}

	if len(names) != 3 {
		t.Errorf("expected 3 results, got %d", len(names))
	}
}

func TestListCompareResults_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	names, err := snapshot.ListCompareResults(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected 0 results, got %d", len(names))
	}
}

func TestListCompareResults_NonExistentDir(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "logsnap-no-such-compare-dir")
	names, err := snapshot.ListCompareResults(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %d", len(names))
	}
}
