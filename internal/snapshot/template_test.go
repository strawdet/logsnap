package snapshot

import (
	"os"
	"testing"
	"time"
)

func TestSaveAndLoadTemplate(t *testing.T) {
	dir := t.TempDir()
	tmpl := &Template{
		Name:   "prod",
		Labels: map[string]string{"env": "production", "team": "platform"},
		Tags:   []string{"release", "stable"},
	}
	if err := SaveTemplate(dir, tmpl); err != nil {
		t.Fatalf("SaveTemplate: %v", err)
	}
	loaded, err := LoadTemplate(dir, "prod")
	if err != nil {
		t.Fatalf("LoadTemplate: %v", err)
	}
	if loaded.Name != "prod" {
		t.Errorf("expected name prod, got %s", loaded.Name)
	}
	if loaded.Labels["env"] != "production" {
		t.Errorf("expected env=production")
	}
	if len(loaded.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(loaded.Tags))
	}
	if loaded.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestLoadTemplate_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadTemplate(dir, "missing")
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestDeleteTemplate(t *testing.T) {
	dir := t.TempDir()
	tmpl := &Template{Name: "staging", Tags: []string{"dev"}}
	_ = SaveTemplate(dir, tmpl)
	if err := DeleteTemplate(dir, "staging"); err != nil {
		t.Fatalf("DeleteTemplate: %v", err)
	}
	if _, err := os.Stat(templatePath(dir, "staging")); !os.IsNotExist(err) {
		t.Error("expected file to be deleted")
	}
}

func TestDeleteTemplate_NotFound(t *testing.T) {
	dir := t.TempDir()
	if err := DeleteTemplate(dir, "ghost"); err == nil {
		t.Fatal("expected error")
	}
}

func TestApplyTemplate(t *testing.T) {
	dir := t.TempDir()
	tmpl := &Template{
		Name:   "base",
		Labels: map[string]string{"env": "staging", "owner": "ops"},
		Tags:   []string{"auto", "nightly"},
	}
	_ = SaveTemplate(dir, tmpl)

	s := &Snapshot{
		ID:        "snap-001",
		CreatedAt: time.Now(),
		Labels:    map[string]string{"env": "prod"},
		Tags:      []string{"auto"},
	}
	if err := ApplyTemplate(dir, "base", s); err != nil {
		t.Fatalf("ApplyTemplate: %v", err)
	}
	if s.Labels["env"] != "prod" {
		t.Error("existing label should not be overwritten")
	}
	if s.Labels["owner"] != "ops" {
		t.Error("expected owner=ops from template")
	}
	if len(s.Tags) != 2 {
		t.Errorf("expected 2 tags (deduped), got %d", len(s.Tags))
	}
}
