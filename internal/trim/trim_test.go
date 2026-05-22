package trim_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/trim"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_HOST":     "db.local",
		"DB_PASSWORD": "secret",
		"DEBUG":       "true",
		"LOG_LEVEL":   "info",
	}
}

func TestApply_ExplicitKeys_RemovesExact(t *testing.T) {
	result := trim.Apply(baseEnv(), trim.Options{Keys: []string{"DEBUG", "LOG_LEVEL"}})
	if _, ok := result.Vars["DEBUG"]; ok {
		t.Error("expected DEBUG to be removed")
	}
	if _, ok := result.Vars["LOG_LEVEL"]; ok {
		t.Error("expected LOG_LEVEL to be removed")
	}
	if len(result.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(result.Removed))
	}
}

func TestApply_PrefixFilter_RemovesMatchingKeys(t *testing.T) {
	result := trim.Apply(baseEnv(), trim.Options{Prefix: "APP_"})
	if _, ok := result.Vars["APP_HOST"]; ok {
		t.Error("expected APP_HOST to be removed")
	}
	if _, ok := result.Vars["APP_PORT"]; ok {
		t.Error("expected APP_PORT to be removed")
	}
	if _, ok := result.Vars["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be preserved")
	}
}

func TestApply_SuffixFilter_RemovesMatchingKeys(t *testing.T) {
	result := trim.Apply(baseEnv(), trim.Options{Suffix: "_HOST"})
	if _, ok := result.Vars["APP_HOST"]; ok {
		t.Error("expected APP_HOST to be removed")
	}
	if _, ok := result.Vars["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed")
	}
	if _, ok := result.Vars["DB_PASSWORD"]; !ok {
		t.Error("expected DB_PASSWORD to be preserved")
	}
}

func TestApply_DryRun_DoesNotMutateMap(t *testing.T) {
	result := trim.Apply(baseEnv(), trim.Options{Prefix: "DB_", DryRun: true})
	if result.Vars != nil {
		t.Error("expected Vars to be nil in dry-run mode")
	}
	if len(result.Removed) == 0 {
		t.Error("expected Removed to be populated in dry-run mode")
	}
}

func TestApply_NoOptions_RemovesNothing(t *testing.T) {
	env := baseEnv()
	result := trim.Apply(env, trim.Options{})
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
	if len(result.Vars) != len(env) {
		t.Errorf("expected %d keys, got %d", len(env), len(result.Vars))
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	trim.Apply(env, trim.Options{Keys: []string{"DEBUG"}})
	if _, ok := env["DEBUG"]; !ok {
		t.Error("Apply must not mutate the original map")
	}
}
