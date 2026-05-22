package clone

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildSummary_Added(t *testing.T) {
	before := map[string]string{}
	after := map[string]string{"NEW_KEY": "val"}
	src := map[string]string{"NEW_KEY": "val"}
	opts := DefaultOptions()

	s := BuildSummary(before, after, src, opts)
	if len(s.Added) != 1 || s.Added[0] != "NEW_KEY" {
		t.Errorf("expected Added=[NEW_KEY], got %v", s.Added)
	}
}

func TestBuildSummary_Skipped(t *testing.T) {
	before := map[string]string{"APP_HOST": "old"}
	after := map[string]string{"APP_HOST": "old"}
	src := map[string]string{"APP_HOST": "new"}
	opts := DefaultOptions()

	s := BuildSummary(before, after, src, opts)
	if len(s.Skipped) != 1 || s.Skipped[0] != "APP_HOST" {
		t.Errorf("expected Skipped=[APP_HOST], got %v", s.Skipped)
	}
}

func TestBuildSummary_Replaced(t *testing.T) {
	before := map[string]string{"APP_HOST": "old"}
	after := map[string]string{"APP_HOST": "new"}
	src := map[string]string{"APP_HOST": "new"}
	opts := DefaultOptions()
	opts.Overwrite = true

	s := BuildSummary(before, after, src, opts)
	if len(s.Replaced) != 1 || s.Replaced[0] != "APP_HOST" {
		t.Errorf("expected Replaced=[APP_HOST], got %v", s.Replaced)
	}
}

func TestWriteText_ContainsPrefixes(t *testing.T) {
	s := Summary{
		Added:    []string{"NEW_KEY"},
		Replaced: []string{"OLD_KEY"},
		Skipped:  []string{"SKIP_KEY"},
	}
	var buf bytes.Buffer
	WriteText(&buf, s)
	out := buf.String()

	if !strings.Contains(out, "+ NEW_KEY") {
		t.Errorf("expected '+ NEW_KEY' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "~ OLD_KEY") {
		t.Errorf("expected '~ OLD_KEY' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "= SKIP_KEY (skipped)") {
		t.Errorf("expected '= SKIP_KEY (skipped)' in output, got:\n%s", out)
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	s := Summary{
		Added:   []string{"A"},
		Skipped: []string{"B"},
	}
	var buf bytes.Buffer
	if err := WriteJSON(&buf, s); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded Summary
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(decoded.Added) != 1 || decoded.Added[0] != "A" {
		t.Errorf("expected Added=[A], got %v", decoded.Added)
	}
}
