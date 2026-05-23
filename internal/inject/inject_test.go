package inject

import (
	"os"
	"testing"
)

func TestInject_SetsNewKeys(t *testing.T) {
	env := map[string]string{"INJECT_TEST_NEW": "hello"}
	os.Unsetenv("INJECT_TEST_NEW")

	results, err := Inject(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Action != "set" {
		t.Fatalf("expected one 'set' result, got %+v", results)
	}
	if got := os.Getenv("INJECT_TEST_NEW"); got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
	os.Unsetenv("INJECT_TEST_NEW")
}

func TestInject_SkipsExistingByDefault(t *testing.T) {
	os.Setenv("INJECT_TEST_EXIST", "original")
	defer os.Unsetenv("INJECT_TEST_EXIST")

	env := map[string]string{"INJECT_TEST_EXIST": "new"}
	results, err := Inject(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Action != "skipped" {
		t.Fatalf("expected skipped, got %+v", results)
	}
	if got := os.Getenv("INJECT_TEST_EXIST"); got != "original" {
		t.Errorf("value should not have changed, got %q", got)
	}
}

func TestInject_OverwriteExisting(t *testing.T) {
	os.Setenv("INJECT_TEST_OW", "old")
	defer os.Unsetenv("INJECT_TEST_OW")

	env := map[string]string{"INJECT_TEST_OW": "new"}
	opts := DefaultOptions()
	opts.Overwrite = true

	results, err := Inject(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Action != "overwritten" {
		t.Fatalf("expected overwritten, got %+v", results)
	}
	if got := os.Getenv("INJECT_TEST_OW"); got != "new" {
		t.Errorf("expected 'new', got %q", got)
	}
}

func TestInject_DryRun_DoesNotMutateEnv(t *testing.T) {
	os.Unsetenv("INJECT_TEST_DRY")

	env := map[string]string{"INJECT_TEST_DRY": "value"}
	opts := DefaultOptions()
	opts.DryRun = true

	results, err := Inject(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Action != "set" {
		t.Fatalf("expected planned set result, got %+v", results)
	}
	if _, exists := os.LookupEnv("INJECT_TEST_DRY"); exists {
		t.Error("dry run should not have set the variable")
		os.Unsetenv("INJECT_TEST_DRY")
	}
}

func TestInject_KeyFilter_OnlySetsListed(t *testing.T) {
	os.Unsetenv("INJECT_A")
	os.Unsetenv("INJECT_B")
	defer func() {
		os.Unsetenv("INJECT_A")
		os.Unsetenv("INJECT_B")
	}()

	env := map[string]string{"INJECT_A": "1", "INJECT_B": "2"}
	opts := DefaultOptions()
	opts.Keys = []string{"INJECT_A"}

	results, err := Inject(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Key != "INJECT_A" {
		t.Fatalf("expected only INJECT_A, got %+v", results)
	}
	if _, exists := os.LookupEnv("INJECT_B"); exists {
		t.Error("INJECT_B should not have been set")
	}
}
