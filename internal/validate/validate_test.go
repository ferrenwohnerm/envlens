package validate_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/validate"
)

func baseEnv() map[string]string {
	return map[string]string{
		"APP_ENV":  "production",
		"PORT":     "8080",
		"LOG_LEVEL": "info",
	}
}

func TestRun_NoFindings_WhenSchemaMatches(t *testing.T) {
	env := baseEnv()
	schema := validate.Schema{
		"APP_ENV":   {Required: true, Pattern: `^(production|staging|development)$`},
		"PORT":      {Required: true, Pattern: `^\d+$`},
		"LOG_LEVEL": {Required: true},
	}
	findings := validate.Run(env, schema)
	for _, f := range findings {
		if f.Severity == "error" {
			t.Errorf("unexpected error finding: %s — %s", f.Key, f.Message)
		}
	}
}

func TestRun_MissingRequiredKey_ProducesError(t *testing.T) {
	env := map[string]string{"PORT": "8080"}
	schema := validate.Schema{
		"APP_ENV": {Required: true},
		"PORT":    {Required: true},
	}
	findings := validate.Run(env, schema)
	found := false
	for _, f := range findings {
		if f.Key == "APP_ENV" && f.Severity == "error" {
			found = true
		}
	}
	if !found {
		t.Error("expected error finding for missing required key APP_ENV")
	}
}

func TestRun_PatternMismatch_ProducesError(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	schema := validate.Schema{
		"PORT": {Required: true, Pattern: `^\d+$`},
	}
	findings := validate.Run(env, schema)
	if len(findings) == 0 {
		t.Fatal("expected at least one finding")
	}
	if findings[0].Severity != "error" {
		t.Errorf("expected error severity, got %q", findings[0].Severity)
	}
}

func TestRun_UndefinedKey_ProducesWarning(t *testing.T) {
	env := map[string]string{"UNKNOWN_KEY": "value"}
	schema := validate.Schema{}
	findings := validate.Run(env, schema)
	if len(findings) == 0 {
		t.Fatal("expected a warning for undefined key")
	}
	if findings[0].Severity != "warning" {
		t.Errorf("expected warning severity, got %q", findings[0].Severity)
	}
}

func TestSummary_NoFindings(t *testing.T) {
	s := validate.Summary(nil)
	if s != "schema validation passed: no findings" {
		t.Errorf("unexpected summary: %q", s)
	}
}

func TestSummary_WithFindings(t *testing.T) {
	findings := []validate.Finding{
		{Key: "A", Message: "missing", Severity: "error"},
		{Key: "B", Message: "undefined", Severity: "warning"},
	}
	s := validate.Summary(findings)
	if s != "schema validation: 1 error(s), 1 warning(s)" {
		t.Errorf("unexpected summary: %q", s)
	}
}
