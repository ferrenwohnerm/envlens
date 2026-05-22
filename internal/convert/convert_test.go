package convert_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envlens/internal/convert"
)

var baseVars = map[string]string{
	"APP_ENV":    "production",
	"DB_HOST":    "localhost",
	"SECRET_KEY": "s3cr3t",
}

func TestConvert_DotenvFormat_ContainsAllKeys(t *testing.T) {
	out, err := convert.Convert(baseVars, convert.Options{TargetFormat: convert.FormatDotenv, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range baseVars {
		if !strings.Contains(out, k) {
			t.Errorf("expected key %q in output", k)
		}
		if !strings.Contains(out, v) {
			t.Errorf("expected value %q in output", v)
		}
	}
}

func TestConvert_DotenvFormat_QuotesSpacedValues(t *testing.T) {
	vars := map[string]string{"GREETING": "hello world"}
	out, err := convert.Convert(vars, convert.Options{TargetFormat: convert.FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value in output, got: %s", out)
	}
}

func TestConvert_JSONFormat_ValidJSON(t *testing.T) {
	out, err := convert.Convert(baseVars, convert.Options{TargetFormat: convert.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]string
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", parsed["APP_ENV"])
	}
}

func TestConvert_ShellFormat_ExportPrefix(t *testing.T) {
	out, err := convert.Convert(baseVars, convert.Options{TargetFormat: convert.FormatShell, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if !strings.HasPrefix(line, "export ") {
			t.Errorf("expected 'export ' prefix, got: %q", line)
		}
	}
}

func TestConvert_YAMLFormat_ContainsKeys(t *testing.T) {
	out, err := convert.Convert(baseVars, convert.Options{TargetFormat: convert.FormatYAML, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k := range baseVars {
		if !strings.Contains(out, k+":") {
			t.Errorf("expected key %q in YAML output", k)
		}
	}
}

func TestConvert_UnknownFormat_ReturnsError(t *testing.T) {
	_, err := convert.Convert(baseVars, convert.Options{TargetFormat: convert.Format("toml")})
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestConvert_SortKeys_ProducesOrderedOutput(t *testing.T) {
	vars := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	out, err := convert.Convert(vars, convert.Options{TargetFormat: convert.FormatDotenv, SortKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line to start with A_KEY, got %q", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected last line to start with Z_KEY, got %q", lines[2])
	}
}
