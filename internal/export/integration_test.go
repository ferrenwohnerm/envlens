package export_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/export"
	"github.com/yourorg/envlens/internal/parser"
)

const sampleEnvFile = `
# database config
DB_HOST=localhost
DB_PORT=5432
DB_NAME=appdb
SECRET="my secret value"
`

func TestExport_RoundTrip_DotenvFormat(t *testing.T) {
	vars, err := parser.ParseFile(strings.NewReader(sampleEnvFile))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	var sb strings.Builder
	opts := export.DefaultOptions()
	if err := export.Write(vars, opts, &sb); err != nil {
		t.Fatalf("export error: %v", err)
	}

	out := sb.String()
	for _, key := range []string{"DB_HOST", "DB_PORT", "DB_NAME", "SECRET"} {
		if !strings.Contains(out, key) {
			t.Errorf("expected key %q in exported output", key)
		}
	}
}

func TestExport_RoundTrip_JSONFormat(t *testing.T) {
	vars, err := parser.ParseFile(strings.NewReader(sampleEnvFile))
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	var sb strings.Builder
	opts := export.Options{Format: export.FormatJSON, Sorted: true}
	if err := export.Write(vars, opts, &sb); err != nil {
		t.Fatalf("export error: %v", err)
	}

	out := sb.String()
	if !strings.Contains(out, `"DB_HOST"`) {
		t.Errorf("expected DB_HOST in JSON output")
	}
	if !strings.Contains(out, `"localhost"`) {
		t.Errorf("expected localhost value in JSON output")
	}
}
