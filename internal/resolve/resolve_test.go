package resolve

import (
	"os"
	"testing"
)

func TestResolve_NoReferences_ReturnsUnchanged(t *testing.T) {
	vars := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	got, err := Resolve(vars, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["HOST"] != "localhost" || got["PORT"] != "5432" {
		t.Errorf("expected unchanged values, got %v", got)
	}
}

func TestResolve_InternalReference_ExpandsValue(t *testing.T) {
	vars := map[string]string{
		"BASE_URL": "http://${HOST}:${PORT}",
		"HOST":     "example.com",
		"PORT":     "8080",
	}
	got, err := Resolve(vars, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["BASE_URL"] != "http://example.com:8080" {
		t.Errorf("expected expanded URL, got %q", got["BASE_URL"])
	}
}

func TestResolve_MissingRef_DefaultOptions_LeavesEmpty(t *testing.T) {
	vars := map[string]string{
		"DSN": "postgres://${DB_USER}@localhost/mydb",
	}
	got, err := Resolve(vars, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DSN"] != "postgres://@localhost/mydb" {
		t.Errorf("unexpected value: %q", got["DSN"])
	}
}

func TestResolve_MissingRef_ErrorOnMissing_ReturnsError(t *testing.T) {
	vars := map[string]string{
		"DSN": "postgres://${DB_USER}@localhost/mydb",
	}
	opts := DefaultOptions()
	opts.ErrorOnMissing = true
	_, err := Resolve(vars, opts)
	if err == nil {
		t.Fatal("expected error for missing reference, got nil")
	}
}

func TestResolve_FallbackToEnv_UsesProcessEnv(t *testing.T) {
	os.Setenv("INJECTED_HOST", "env-host")
	defer os.Unsetenv("INJECTED_HOST")

	vars := map[string]string{
		"ADDR": "${INJECTED_HOST}:9000",
	}
	opts := DefaultOptions()
	opts.FallbackToEnv = true

	got, err := Resolve(vars, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["ADDR"] != "env-host:9000" {
		t.Errorf("expected env fallback, got %q", got["ADDR"])
	}
}

func TestResolve_EmptyMap_ReturnsEmpty(t *testing.T) {
	got, err := Resolve(map[string]string{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
