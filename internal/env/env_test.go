package env_test

import (
	"os"
	"testing"

	"github.com/yourusername/envlens/internal/env"
)

func TestInject_SetsVariables(t *testing.T) {
	t.Setenv("ENVLENS_TEST_INJECT", "") // ensure cleanup
	os.Unsetenv("ENVLENS_TEST_INJECT")

	vars := map[string]string{"ENVLENS_TEST_INJECT": "hello"}
	if err := env.Inject(vars, env.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("ENVLENS_TEST_INJECT"); got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestInject_SkipsExistingByDefault(t *testing.T) {
	t.Setenv("ENVLENS_EXISTING", "original")

	vars := map[string]string{"ENVLENS_EXISTING": "replaced"}
	if err := env.Inject(vars, env.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("ENVLENS_EXISTING"); got != "original" {
		t.Errorf("expected %q, got %q", "original", got)
	}
}

func TestInject_Overwrite_ReplacesExisting(t *testing.T) {
	t.Setenv("ENVLENS_OW", "old")

	opts := env.DefaultOptions()
	opts.Overwrite = true
	vars := map[string]string{"ENVLENS_OW": "new"}
	if err := env.Inject(vars, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("ENVLENS_OW"); got != "new" {
		t.Errorf("expected %q, got %q", "new", got)
	}
}

func TestExtract_ReturnsRequestedKeys(t *testing.T) {
	t.Setenv("ENVLENS_EX_A", "alpha")
	t.Setenv("ENVLENS_EX_B", "beta")

	got, err := env.Extract([]string{"ENVLENS_EX_A", "ENVLENS_EX_B"}, env.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["ENVLENS_EX_A"] != "alpha" || got["ENVLENS_EX_B"] != "beta" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestExtract_MissingKey_FailOnMissing_ReturnsError(t *testing.T) {
	os.Unsetenv("ENVLENS_MISSING")
	opts := env.DefaultOptions()
	opts.FailOnMissing = true
	_, err := env.Extract([]string{"ENVLENS_MISSING"}, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestExtract_MissingKey_DefaultOptions_OmitsKey(t *testing.T) {
	os.Unsetenv("ENVLENS_ABSENT")
	got, err := env.Extract([]string{"ENVLENS_ABSENT"}, env.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["ENVLENS_ABSENT"]; ok {
		t.Error("expected key to be absent from result")
	}
}

func TestSnapshot_CapturesPrefix(t *testing.T) {
	t.Setenv("ENVLENS_SNAP_X", "x")
	t.Setenv("ENVLENS_SNAP_Y", "y")

	snap := env.Snapshot("ENVLENS_SNAP_")
	if snap["ENVLENS_SNAP_X"] != "x" || snap["ENVLENS_SNAP_Y"] != "y" {
		t.Errorf("unexpected snapshot: %v", snap)
	}
}

func TestSortedKeys_ReturnsAlphabeticalOrder(t *testing.T) {
	vars := map[string]string{"Z": "1", "A": "2", "M": "3"}
	keys := env.SortedKeys(vars)
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], k)
		}
	}
}
