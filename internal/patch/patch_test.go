package patch_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/patch"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "staging",
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"LOG_LEVEL": "debug",
	}
}

func TestApply_SetNewKey_AddsEntry(t *testing.T) {
	out, err := patch.Apply(baseEnv(), []patch.Instruction{
		{Op: patch.OpSet, Key: "NEW_KEY", Value: "hello"},
	}, patch.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["NEW_KEY"]; got != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", got)
	}
}

func TestApply_SetExistingKey_OverwritesValue(t *testing.T) {
	out, err := patch.Apply(baseEnv(), []patch.Instruction{
		{Op: patch.OpSet, Key: "APP_ENV", Value: "production"},
	}, patch.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["APP_ENV"]; got != "production" {
		t.Errorf("expected APP_ENV=production, got %q", got)
	}
}

func TestApply_UnsetKey_RemovesEntry(t *testing.T) {
	out, err := patch.Apply(baseEnv(), []patch.Instruction{
		{Op: patch.OpUnset, Key: "LOG_LEVEL"},
	}, patch.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["LOG_LEVEL"]; ok {
		t.Error("expected LOG_LEVEL to be removed")
	}
}

func TestApply_UnsetMissingKey_ErrorOnMissing(t *testing.T) {
	opts := patch.DefaultOptions()
	opts.ErrorOnMissing = true
	_, err := patch.Apply(baseEnv(), []patch.Instruction{
		{Op: patch.OpUnset, Key: "DOES_NOT_EXIST"},
	}, opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestApply_RenameKey_MovesValue(t *testing.T) {
	out, err := patch.Apply(baseEnv(), []patch.Instruction{
		{Op: patch.OpRename, Key: "DB_HOST", To: "DATABASE_HOST"},
	}, patch.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("old key DB_HOST should have been removed")
	}
	if got := out["DATABASE_HOST"]; got != "localhost" {
		t.Errorf("expected DATABASE_HOST=localhost, got %q", got)
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	_, err := patch.Apply(env, []patch.Instruction{
		{Op: patch.OpSet, Key: "APP_ENV", Value: "production"},
		{Op: patch.OpUnset, Key: "LOG_LEVEL"},
	}, patch.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["APP_ENV"] != "staging" {
		t.Error("original map was mutated: APP_ENV changed")
	}
	if _, ok := env["LOG_LEVEL"]; !ok {
		t.Error("original map was mutated: LOG_LEVEL removed")
	}
}

func TestApply_UnknownOp_ReturnsError(t *testing.T) {
	_, err := patch.Apply(baseEnv(), []patch.Instruction{
		{Op: patch.Op("explode"), Key: "APP_ENV"},
	}, patch.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for unknown op, got nil")
	}
}
