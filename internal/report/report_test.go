package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envlens/internal/diff"
	"github.com/user/envlens/internal/report"
)

func makeResult() diff.Result {
	return diff.Result{
		"APP_ENV":  {Status: diff.Identical, ValueA: "production", ValueB: "production"},
		"DB_HOST":  {Status: diff.Changed, ValueA: "localhost", ValueB: "db.prod.internal"},
		"API_KEY":  {Status: diff.OnlyInA, ValueA: "abc123", ValueB: ""},
		"NEW_FLAG": {Status: diff.OnlyInB, ValueA: "", ValueB: "true"},
	}
}

func TestWrite_TextFormat_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	err := report.Write(&buf, makeResult(), report.FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, expected := range []string{"APP_ENV", "DB_HOST", "API_KEY", "NEW_FLAG"} {
		if !strings.Contains(out, expected) {
			t.Errorf("expected output to contain %q, got:\n%s", expected, out)
		}
	}
}

func TestWrite_TextFormat_StatusPrefixes(t *testing.T) {
	var buf bytes.Buffer
	_ = report.Write(&buf, makeResult(), report.FormatText)
	out := buf.String()

	if !strings.Contains(out, "- API_KEY") {
		t.Errorf("expected '- API_KEY' in output")
	}
	if !strings.Contains(out, "+ NEW_FLAG") {
		t.Errorf("expected '+ NEW_FLAG' in output")
	}
	if !strings.Contains(out, "~ DB_HOST") {
		t.Errorf("expected '~ DB_HOST' in output")
	}
}

func TestWrite_JSONFormat_ContainsKeys(t *testing.T) {
	var buf bytes.Buffer
	err := report.Write(&buf, makeResult(), report.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, expected := range []string{"APP_ENV", "DB_HOST", "API_KEY", "NEW_FLAG"} {
		if !strings.Contains(out, expected) {
			t.Errorf("expected JSON output to contain %q, got:\n%s", expected, out)
		}
	}
}

func TestWrite_JSONFormat_StartsAndEndsWithBraces(t *testing.T) {
	var buf bytes.Buffer
	_ = report.Write(&buf, makeResult(), report.FormatJSON)
	out := strings.TrimSpace(buf.String())
	if !strings.HasPrefix(out, "{") || !strings.HasSuffix(out, "}") {
		t.Errorf("expected JSON output wrapped in braces, got:\n%s", out)
	}
}

func TestWrite_EmptyResult(t *testing.T) {
	var buf bytes.Buffer
	err := report.Write(&buf, diff.Result{}, report.FormatText)
	if err != nil {
		t.Fatalf("unexpected error on empty result: %v", err)
	}
}
