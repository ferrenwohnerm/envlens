package split

import (
	"testing"
)

var baseEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "envlens",
	"APP_VERSION": "1.0.0",
	"SECRET_KEY":  "abc123",
}

func TestRun_SplitsIntoGroups(t *testing.T) {
	result := Run(baseEnv, []string{"DB_", "APP_"}, DefaultOptions())

	if len(result.Groups["DB_"]) != 2 {
		t.Errorf("expected 2 DB_ keys, got %d", len(result.Groups["DB_"]))
	}
	if len(result.Groups["APP_"]) != 2 {
		t.Errorf("expected 2 APP_ keys, got %d", len(result.Groups["APP_"]))
	}
}

func TestRun_StripPrefix_RemovesPrefixFromKey(t *testing.T) {
	opts := DefaultOptions()
	opts.StripPrefix = true
	result := Run(baseEnv, []string{"DB_"}, opts)

	if _, ok := result.Groups["DB_"]["HOST"]; !ok {
		t.Error("expected key HOST after stripping DB_ prefix")
	}
	if _, ok := result.Groups["DB_"]["PORT"]; !ok {
		t.Error("expected key PORT after stripping DB_ prefix")
	}
}

func TestRun_NoStripPrefix_KeepsFullKey(t *testing.T) {
	opts := DefaultOptions()
	opts.StripPrefix = false
	result := Run(baseEnv, []string{"DB_"}, opts)

	if _, ok := result.Groups["DB_"]["DB_HOST"]; !ok {
		t.Error("expected full key DB_HOST when StripPrefix=false")
	}
}

func TestRun_UnmatchedKeys_NotInGroupsByDefault(t *testing.T) {
	result := Run(baseEnv, []string{"DB_", "APP_"}, DefaultOptions())

	if len(result.Unmatched) != 1 {
		t.Errorf("expected 1 unmatched key, got %d", len(result.Unmatched))
	}
	if result.Unmatched[0] != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY unmatched, got %s", result.Unmatched[0])
	}
	if _, ok := result.Groups[""]; ok {
		t.Error("expected no empty-string group when IncludeUnmatched=false")
	}
}

func TestRun_IncludeUnmatched_AddsEmptyGroup(t *testing.T) {
	opts := DefaultOptions()
	opts.IncludeUnmatched = true
	result := Run(baseEnv, []string{"DB_", "APP_"}, opts)

	if result.Groups[""] == nil {
		t.Fatal("expected empty-string group for unmatched keys")
	}
	if _, ok := result.Groups[""]["SECRET_KEY"]; !ok {
		t.Error("expected SECRET_KEY in unmatched group")
	}
}

func TestRun_FirstPrefixWins_OnOverlap(t *testing.T) {
	vars := map[string]string{
		"APP_SERVICE_NAME": "svc",
	}
	opts := DefaultOptions()
	opts.StripPrefix = false
	result := Run(vars, []string{"APP_", "APP_SERVICE_"}, opts)

	if _, ok := result.Groups["APP_"]["APP_SERVICE_NAME"]; !ok {
		t.Error("expected APP_ group to win over APP_SERVICE_ (first match wins)")
	}
	if len(result.Groups["APP_SERVICE_"]) != 0 {
		t.Error("expected APP_SERVICE_ group to be empty")
	}
}

func TestRun_EmptyPrefixes_AllUnmatched(t *testing.T) {
	opts := DefaultOptions()
	opts.IncludeUnmatched = true
	result := Run(baseEnv, []string{}, opts)

	if len(result.Unmatched) != len(baseEnv) {
		t.Errorf("expected all %d keys unmatched, got %d", len(baseEnv), len(result.Unmatched))
	}
}
