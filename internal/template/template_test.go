package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func baseVars() map[string]string {
	return map[string]string{
		"APP_NAME":     "myapp",
		"APP_PORT":     "8080",
		"DB_PASSWORD":  "supersecret",
		"API_KEY":      "abc123",
		"LOG_LEVEL":    "info",
	}
}

func TestGenerate_DefaultOptions_ReplacesAllWithPlaceholder(t *testing.T) {
	out := Generate(baseVars(), DefaultOptions())
	if strings.Contains(out, "supersecret") {
		t.Error("expected sensitive value to be replaced")
	}
	if strings.Contains(out, "abc123") {
		t.Error("expected API_KEY value to be replaced")
	}
	if !strings.Contains(out, "CHANGEME") {
		t.Error("expected placeholder CHANGEME in output")
	}
}

func TestGenerate_IncludeValues_PreservesNonSensitive(t *testing.T) {
	opts := Options{IncludeValues: true, Placeholder: "CHANGEME"}
	out := Generate(baseVars(), opts)
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Error("expected non-sensitive value to be preserved")
	}
	if strings.Contains(out, "supersecret") {
		t.Error("expected DB_PASSWORD to remain redacted")
	}
}

func TestGenerate_SensitiveKeys_AlwaysRedacted(t *testing.T) {
	opts := Options{IncludeValues: true, Placeholder: "CHANGEME"}
	out := Generate(baseVars(), opts)
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "DB_PASSWORD") && !strings.Contains(line, "CHANGEME") {
			t.Errorf("DB_PASSWORD should be redacted, got: %s", line)
		}
		if strings.HasPrefix(line, "API_KEY") && !strings.Contains(line, "CHANGEME") {
			t.Errorf("API_KEY should be redacted, got: %s", line)
		}
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	opts := Options{Placeholder: "TODO"}
	out := Generate(map[string]string{"APP_NAME": "x"}, opts)
	if !strings.Contains(out, "TODO") {
		t.Error("expected custom placeholder TODO")
	}
}

func TestGenerate_OutputIsSorted(t *testing.T) {
	out := Generate(baseVars(), DefaultOptions())
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i := 1; i < len(lines); i++ {
		if lines[i] < lines[i-1] {
			t.Errorf("output not sorted: %q before %q", lines[i-1], lines[i])
		}
	}
}

func TestWriteFile_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env.template")
	err := WriteFile(path, baseVars(), DefaultOptions())
	if err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), "APP_NAME") {
		t.Error("expected APP_NAME in written file")
	}
}
