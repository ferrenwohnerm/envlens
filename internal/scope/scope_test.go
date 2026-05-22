package scope_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/scope"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"DB_HOST":     "db.internal",
		"DB_PASSWORD": "secret",
		"LOG_LEVEL":   "info",
	}
}

func TestRun_NoPrefixes_ReturnsAll(t *testing.T) {
	result := scope.Run(baseEnv(), scope.DefaultOptions())
	if len(result.Vars) != 5 {
		t.Fatalf("expected 5 vars, got %d", len(result.Vars))
	}
	if len(result.Dropped) != 0 {
		t.Fatalf("expected 0 dropped, got %d", len(result.Dropped))
	}
}

func TestRun_PrefixFilter_RetainsMatchingKeys(t *testing.T) {
	opts := scope.Options{Prefixes: []string{"APP_"}}
	result := scope.Run(baseEnv(), opts)
	if len(result.Vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(result.Vars))
	}
	if _, ok := result.Vars["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := result.Vars["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in result")
	}
}

func TestRun_PrefixFilter_DropsNonMatchingKeys(t *testing.T) {
	opts := scope.Options{Prefixes: []string{"APP_"}}
	result := scope.Run(baseEnv(), opts)
	if len(result.Dropped) != 3 {
		t.Fatalf("expected 3 dropped keys, got %d", len(result.Dropped))
	}
}

func TestRun_StripPrefix_RemovesPrefixFromKey(t *testing.T) {
	opts := scope.Options{Prefixes: []string{"APP_"}, StripPrefix: true}
	result := scope.Run(baseEnv(), opts)
	if _, ok := result.Vars["HOST"]; !ok {
		t.Error("expected HOST after stripping APP_ prefix")
	}
	if _, ok := result.Vars["PORT"]; !ok {
		t.Error("expected PORT after stripping APP_ prefix")
	}
}

func TestRun_MultiplePrefixes_RetainsAll(t *testing.T) {
	opts := scope.Options{Prefixes: []string{"APP_", "DB_"}}
	result := scope.Run(baseEnv(), opts)
	if len(result.Vars) != 4 {
		t.Fatalf("expected 4 vars, got %d", len(result.Vars))
	}
}

func TestRun_CaseFold_MatchesCaseInsensitive(t *testing.T) {
	env := map[string]string{"app_host": "localhost", "DB_HOST": "db"}
	opts := scope.Options{Prefixes: []string{"APP_"}, CaseFold: true}
	result := scope.Run(env, opts)
	if len(result.Vars) != 1 {
		t.Fatalf("expected 1 var, got %d", len(result.Vars))
	}
	if _, ok := result.Vars["app_host"]; !ok {
		t.Error("expected app_host in result")
	}
}

func TestRun_EmptyEnv_ReturnsEmpty(t *testing.T) {
	opts := scope.Options{Prefixes: []string{"APP_"}}
	result := scope.Run(map[string]string{}, opts)
	if len(result.Vars) != 0 {
		t.Fatalf("expected 0 vars, got %d", len(result.Vars))
	}
}
