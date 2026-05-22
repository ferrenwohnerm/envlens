package rename

import (
	"testing"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_HOST":     "db.internal",
		"DB_PASSWORD": "secret",
		"DEBUG":       "true",
	}
}

func TestRename_ExplicitMapping_RenamesKeys(t *testing.T) {
	env := baseEnv()
	opts := Options{Mapping: map[string]string{"APP_HOST": "SERVICE_HOST"}}

	out, results, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("old key APP_HOST should have been removed")
	}
	if out["SERVICE_HOST"] != "localhost" {
		t.Errorf("expected SERVICE_HOST=localhost, got %q", out["SERVICE_HOST"])
	}
	if len(results) != 1 || results[0].OldKey != "APP_HOST" || results[0].NewKey != "SERVICE_HOST" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestRename_PrefixSubstitution_RenamesMatchingKeys(t *testing.T) {
	env := baseEnv()
	opts := Options{OldPrefix: "APP_", NewPrefix: "SVC_"}

	out, results, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("APP_HOST should have been renamed")
	}
	if out["SVC_HOST"] != "localhost" {
		t.Errorf("expected SVC_HOST=localhost, got %q", out["SVC_HOST"])
	}
	if out["SVC_PORT"] != "8080" {
		t.Errorf("expected SVC_PORT=8080, got %q", out["SVC_PORT"])
	}
	if len(results) != 2 {
		t.Errorf("expected 2 renames, got %d", len(results))
	}
}

func TestRename_MissingKey_ErrorOnMissing(t *testing.T) {
	env := baseEnv()
	opts := Options{
		Mapping:        map[string]string{"NONEXISTENT": "NEW_KEY"},
		ErrorOnMissing: true,
	}

	_, _, err := Rename(env, opts)
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestRename_MissingKey_NoErrorByDefault(t *testing.T) {
	env := baseEnv()
	opts := Options{Mapping: map[string]string{"NONEXISTENT": "NEW_KEY"}}

	_, results, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestRename_OriginalMapUnmodified(t *testing.T) {
	env := baseEnv()
	opts := Options{Mapping: map[string]string{"DEBUG": "VERBOSE"}}

	_, _, _ = Rename(env, opts)
	if _, ok := env["DEBUG"]; !ok {
		t.Error("original map should not be modified")
	}
}

func TestRename_SameKeyMapping_NoOpAndNoResult(t *testing.T) {
	env := baseEnv()
	opts := Options{Mapping: map[string]string{"DEBUG": "DEBUG"}}

	_, results, err := Rename(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for no-op rename, got %d", len(results))
	}
}
