package snapshot

import (
	"os"
	"testing"
	"time"
)

func writeTriggerSnapshot(t *testing.T, dir string) *Snapshot {
	t.Helper()
	snap := &Snapshot{
		ID:        "snap-trigger-001",
		CreatedAt: time.Now().UTC(),
		Entries: []LogEntry{
			{Level: "error", Message: "disk full"},
			{Level: "error", Message: "OOM"},
			{Level: "info", Message: "started"},
			{Level: "info", Message: "ok"},
		},
	}
	if err := snap.Save(dir); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	return snap
}

func TestSaveTrigger_AndLoad(t *testing.T) {
	dir := t.TempDir()
	cond := TriggerCondition{ErrorRateThreshold: 0.3, MinEntries: 2, Label: "ci"}
	tr, err := SaveTrigger(dir, "high-error", cond)
	if err != nil {
		t.Fatalf("SaveTrigger: %v", err)
	}
	if tr.Name != "high-error" {
		t.Errorf("expected name high-error, got %s", tr.Name)
	}
	loaded, err := LoadTrigger(dir, "high-error")
	if err != nil {
		t.Fatalf("LoadTrigger: %v", err)
	}
	if loaded.Condition.ErrorRateThreshold != 0.3 {
		t.Errorf("threshold mismatch: got %v", loaded.Condition.ErrorRateThreshold)
	}
	if loaded.Condition.Label != "ci" {
		t.Errorf("label mismatch: got %s", loaded.Condition.Label)
	}
}

func TestLoadTrigger_NotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadTrigger(dir, "ghost")
	if err == nil {
		t.Fatal("expected error for missing trigger")
	}
}

func TestDeleteTrigger(t *testing.T) {
	dir := t.TempDir()
	_, err := SaveTrigger(dir, "tmp", TriggerCondition{MinEntries: 1})
	if err != nil {
		t.Fatalf("SaveTrigger: %v", err)
	}
	if err := DeleteTrigger(dir, "tmp"); err != nil {
		t.Fatalf("DeleteTrigger: %v", err)
	}
	_, err = LoadTrigger(dir, "tmp")
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}

func TestDeleteTrigger_NotFound(t *testing.T) {
	dir := t.TempDir()
	if err := DeleteTrigger(dir, "nope"); err == nil {
		t.Fatal("expected error for missing trigger")
	}
}

func TestEvaluateTrigger_ErrorRate(t *testing.T) {
	dir := t.TempDir()
	snap := writeTriggerSnapshot(t, dir)
	// 2/4 = 0.5 error rate
	cond := TriggerCondition{ErrorRateThreshold: 0.4}
	if !EvaluateTrigger(snap, cond) {
		t.Error("expected trigger to fire")
	}
	cond2 := TriggerCondition{ErrorRateThreshold: 0.6}
	if EvaluateTrigger(snap, cond2) {
		t.Error("expected trigger not to fire")
	}
}

func TestEvaluateTrigger_MinEntries(t *testing.T) {
	dir := t.TempDir()
	snap := writeTriggerSnapshot(t, dir)
	if !EvaluateTrigger(snap, TriggerCondition{MinEntries: 4}) {
		t.Error("expected trigger to fire with MinEntries=4")
	}
	if EvaluateTrigger(snap, TriggerCondition{MinEntries: 10}) {
		t.Error("expected trigger not to fire with MinEntries=10")
	}
}

func TestEvaluateTrigger_NilSnapshot(t *testing.T) {
	if EvaluateTrigger(nil, TriggerCondition{MinEntries: 1}) {
		t.Error("expected false for nil snapshot")
	}
	_ = os.Getenv("CI") // suppress unused import
}
