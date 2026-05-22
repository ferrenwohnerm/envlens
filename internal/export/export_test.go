package export

import (
	"strings"
	"testing"
)

var sampleVars = map[string]string{
	"APP_ENV":    "production",
	"DB_URL":     "postgres://localhost/mydb",
	"SECRET_KEY": "p@ss w0rd",
}

func TestWrite_DotenvFormat_ContainsAllKeys(t *testing.T) {
	var sb strings.Builder
	err := Write(sampleVars, DefaultOptions(), &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, key := range []string{"APP_ENV", "DB_URL", "SECRET_KEY"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in output", key)
		}
	}
}

func TestWrite_DotenvFormat_QuotesValueWithSpaces(t *testing.T) {
	var sb strings.Builder
	_ = Write(map[string]string{"KEY": "hello world"}, DefaultOptions(), &sb)
	if !strings.Contains(sb.String(), `"hello world"`) {
		t.Errorf("expected quoted value in output, got: %s", sb.String())
	}
}

func TestWrite_JSONFormat_ValidJSON(t *testing.T) {
	var sb strings.Builder
	opts := Options{Format: FormatJSON, Sorted: true}
	err := Write(sampleVars, opts, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON object, got: %s", out)
	}
	if !strings.Contains(out, `"APP_ENV"`) {
		t.Errorf("expected APP_ENV in JSON output")
	}
}

func TestWrite_ShellFormat_ExportPrefix(t *testing.T) {
	var sb strings.Builder
	opts := Options{Format: FormatShell, Sorted: true}
	err := Write(sampleVars, opts, &sb)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	for _, l := range lines {
		if !strings.HasPrefix(l, "export ") {
			t.Errorf("expected 'export ' prefix, got: %s", l)
		}
	}
}

func TestWrite_UnknownFormat_ReturnsError(t *testing.T) {
	var sb strings.Builder
	err := Write(sampleVars, Options{Format: "xml"}, &sb)
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestWrite_SortedOutput_IsAlphabetical(t *testing.T) {
	var sb strings.Builder
	_ = Write(sampleVars, DefaultOptions(), &sb)
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	if !strings.HasPrefix(lines[0], "APP_ENV") {
		t.Errorf("expected APP_ENV first in sorted output, got: %s", lines[0])
	}
}
