package scope_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yourorg/envlens/internal/scope"
)

func makeResult() scope.Result {
	return scope.Result{
		Vars:    map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"},
		Dropped: []string{"DB_HOST", "LOG_LEVEL"},
	}
}

func TestWriteText_ContainsScopedVars(t *testing.T) {
	var buf bytes.Buffer
	if err := scope.WriteText(&buf, makeResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "APP_HOST") {
		t.Error("expected APP_HOST in text output")
	}
	if !strings.Contains(out, "APP_PORT") {
		t.Error("expected APP_PORT in text output")
	}
}

func TestWriteText_ContainsDroppedKeys(t *testing.T) {
	var buf bytes.Buffer
	scope.WriteText(&buf, makeResult())
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in dropped section")
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Error("expected LOG_LEVEL in dropped section")
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := scope.WriteJSON(&buf, makeResult()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestWriteJSON_ContainsVarsAndDropped(t *testing.T) {
	var buf bytes.Buffer
	scope.WriteJSON(&buf, makeResult())
	out := buf.String()
	if !strings.Contains(out, "\"vars\"") {
		t.Error("expected 'vars' key in JSON output")
	}
	if !strings.Contains(out, "\"dropped\"") {
		t.Error("expected 'dropped' key in JSON output")
	}
}

func TestWriteJSON_EmptyResult_StillValid(t *testing.T) {
	var buf bytes.Buffer
	r := scope.Result{Vars: map[string]string{}, Dropped: nil}
	if err := scope.WriteJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}
