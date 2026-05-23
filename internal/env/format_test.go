package env_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourusername/envlens/internal/env"
)

func TestWriteText_ContainsKeyValuePairs(t *testing.T) {
	vars := map[string]string{"FOO": "bar", "BAZ": "qux"}
	var buf bytes.Buffer
	if err := env.WriteText(&buf, vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got:\n%s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in output, got:\n%s", out)
	}
}

func TestWriteText_OutputIsSorted(t *testing.T) {
	vars := map[string]string{"Z_KEY": "z", "A_KEY": "a"}
	var buf bytes.Buffer
	if err := env.WriteText(&buf, vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected A_KEY first, got %q", lines[0])
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	vars := map[string]string{"HOST": "localhost", "PORT": "5432"}
	var buf bytes.Buffer
	if err := env.WriteJSON(&buf, vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed []map[string]string
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, buf.String())
	}
	if len(parsed) != 2 {
		t.Errorf("expected 2 entries, got %d", len(parsed))
	}
}

func TestWriteJSON_ContainsExpectedKeys(t *testing.T) {
	vars := map[string]string{"DB_URL": "postgres://localhost/mydb"}
	var buf bytes.Buffer
	if err := env.WriteJSON(&buf, vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "DB_URL") {
		t.Errorf("expected DB_URL in JSON output, got:\n%s", buf.String())
	}
}

func TestWriteSummary_ContainsCount(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	var buf bytes.Buffer
	if err := env.WriteSummary(&buf, vars); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "3 variable(s)") {
		t.Errorf("expected count in summary, got: %s", buf.String())
	}
}
