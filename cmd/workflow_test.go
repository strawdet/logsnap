package cmd

import (
	"testing"
)

func TestParseStep_NoParams(t *testing.T) {
	step, err := parseStep("diff")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if step.Action != "diff" {
		t.Errorf("action: got %q, want %q", step.Action, "diff")
	}
	if step.Name != "diff" {
		t.Errorf("name: got %q, want %q", step.Name, "diff")
	}
	if len(step.Params) != 0 {
		t.Errorf("params: expected empty, got %v", step.Params)
	}
}

func TestParseStep_WithParams(t *testing.T) {
	step, err := parseStep("export:format=json,dest=/tmp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if step.Action != "export" {
		t.Errorf("action: got %q", step.Action)
	}
	if step.Params["format"] != "json" {
		t.Errorf("format: got %q", step.Params["format"])
	}
	if step.Params["dest"] != "/tmp" {
		t.Errorf("dest: got %q", step.Params["dest"])
	}
}

func TestParseStep_InvalidParam(t *testing.T) {
	_, err := parseStep("capture:badparam")
	if err == nil {
		t.Fatal("expected error for invalid param, got nil")
	}
}

func TestParseStep_SingleParam(t *testing.T) {
	step, err := parseStep("capture:label=v1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if step.Params["label"] != "v1" {
		t.Errorf("label: got %q", step.Params["label"])
	}
}
