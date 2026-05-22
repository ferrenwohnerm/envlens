package template_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/parser"
	"github.com/yourorg/envlens/internal/template"
)

func TestGenerate_RoundTrip_ParsedFile(t *testing.T) {
	const raw = `
APP_NAME=myapp
APP_PORT=9000
DB_PASSWORD=secret123
LOG_LEVEL=debug
`
	tmpFile := writeTempFile(t, raw)
	vars, err := parser.ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}

	opts := template.Options{IncludeValues: false, Placeholder: "CHANGEME"}
	out := template.Generate(vars, opts)

	if strings.Contains(out, "secret123") {
		t.Error("expected sensitive value to be stripped from template")
	}
	if !strings.Contains(out, "APP_NAME=CHANGEME") {
		t.Error("expected APP_NAME stub in output")
	}
	if !strings.Contains(out, "DB_PASSWORD=CHANGEME") {
		t.Error("expected DB_PASSWORD stub in output")
	}
}

func TestGenerate_RoundTrip_WithValues(t *testing.T) {
	const raw = `
SERVICE=api
VERSION=1.2.3
AUTH_TOKEN=tok_xyz
`
	tmpFile := writeTempFile(t, raw)
	vars, err := parser.ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile error: %v", err)
	}

	opts := template.Options{IncludeValues: true, Placeholder: "CHANGEME"}
	out := template.Generate(vars, opts)

	if !strings.Contains(out, "SERVICE=api") {
		t.Error("expected SERVICE value preserved")
	}
	if strings.Contains(out, "tok_xyz") {
		t.Error("expected AUTH_TOKEN to remain redacted")
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := t.TempDir(), error(nil)
	_ = err
	path := f + "/.env"
	if writeErr := writeFile(path, content); writeErr != nil {
		t.Fatalf("failed to write temp file: %v", writeErr)
	}
	return path
}

func writeFile(path, content string) error {
	import_os_WriteFile := func(p, c string) error {
		import "os"
		return os.WriteFile(p, []byte(c), 0o644)
	}
	return import_os_WriteFile(path, content)
}
