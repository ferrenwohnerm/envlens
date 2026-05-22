package redact_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/redact"
)

var baseEnv = map[string]string{
	"APP_NAME":      "myapp",
	"DB_PASSWORD":   "s3cr3t",
	"API_KEY":       "abc123",
	"DATABASE_URL":  "postgres://user:pass@host/db",
	"PORT":          "8080",
	"AUTH_TOKEN":    "tok_xyz",
}

func TestApply_DefaultOptions_RedactsSensitiveKeys(t *testing.T) {
	out := redact.Apply(baseEnv, redact.DefaultOptions())
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected API_KEY to be redacted, got %q", out["API_KEY"])
	}
	if out["AUTH_TOKEN"] != "[REDACTED]" {
		t.Errorf("expected AUTH_TOKEN to be redacted, got %q", out["AUTH_TOKEN"])
	}
}

func TestApply_DefaultOptions_PreservesNonSensitiveKeys(t *testing.T) {
	out := redact.Apply(baseEnv, redact.DefaultOptions())
	if out["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be preserved, got %q", out["APP_NAME"])
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT to be preserved, got %q", out["PORT"])
	}
}

func TestApply_DefaultOptions_DSNRule(t *testing.T) {
	out := redact.Apply(baseEnv, redact.DefaultOptions())
	if out["DATABASE_URL"] != "[REDACTED-DSN]" {
		t.Errorf("expected DATABASE_URL to use DSN placeholder, got %q", out["DATABASE_URL"])
	}
}

func TestApply_MaskAll_RedactsEverything(t *testing.T) {
	opts := redact.Options{MaskAll: true, DefaultReplacement: "***"}
	out := redact.Apply(baseEnv, opts)
	for k, v := range out {
		if v != "***" {
			t.Errorf("MaskAll: expected key %q to be \"***\", got %q", k, v)
		}
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	input := map[string]string{"SECRET": "original"}
	_ = redact.Apply(input, redact.DefaultOptions())
	if input["SECRET"] != "original" {
		t.Error("Apply must not mutate the input map")
	}
}

func TestApply_EmptyMap_ReturnsEmpty(t *testing.T) {
	out := redact.Apply(map[string]string{}, redact.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty map, got %d entries", len(out))
	}
}
