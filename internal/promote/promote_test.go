package promote

import (
	"testing"
)

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func TestPromote_AddsMissingKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "old"}

	res, err := Promote(src, dst, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst["B"] != "2" {
		t.Errorf("expected B=2 in dst, got %q", dst["B"])
	}

	actions := actionMap(res)
	if actions["B"] != "add" {
		t.Errorf("expected action add for B, got %q", actions["B"])
	}
}

func TestPromote_SkipsExistingByDefault(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	_, err := Promote(src, dst, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst["A"] != "old" {
		t.Errorf("expected A to remain 'old', got %q", dst["A"])
	}
}

func TestPromote_OverwriteExisting_ReplacesValue(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	opts := DefaultOptions()
	opts.OverwriteExisting = true

	res, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dst["A"] != "new" {
		t.Errorf("expected A=new, got %q", dst["A"])
	}
	actions := actionMap(res)
	if actions["A"] != "overwrite" {
		t.Errorf("expected action overwrite, got %q", actions["A"])
	}
}

func TestPromote_DryRun_DoesNotModifyDst(t *testing.T) {
	src := map[string]string{"NEW": "value"}
	dst := map[string]string{}
	opts := DefaultOptions()
	opts.DryRun = true

	res, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := dst["NEW"]; ok {
		t.Error("dry run must not modify dst")
	}
	actions := actionMap(res)
	if actions["NEW"] != "add" {
		t.Errorf("expected action add in dry-run result, got %q", actions["NEW"])
	}
}

func TestPromote_SkipKeys_ExcludesKey(t *testing.T) {
	src := map[string]string{"SECRET": "s3cr3t", "HOST": "prod"}
	dst := map[string]string{}
	opts := DefaultOptions()
	opts.SkipKeys = map[string]bool{"SECRET": true}

	_, err := Promote(src, dst, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := dst["SECRET"]; ok {
		t.Error("SECRET should have been skipped")
	}
	if dst["HOST"] != "prod" {
		t.Errorf("expected HOST=prod, got %q", dst["HOST"])
	}
}

func TestPromote_NilSrc_ReturnsError(t *testing.T) {
	_, err := Promote(nil, map[string]string{}, DefaultOptions())
	if err == nil {
		t.Error("expected error for nil src")
	}
}

// actionMap builds a key→action lookup from a Result for easy assertions.
func actionMap(r Result) map[string]string {
	m := make(map[string]string, len(r.Changes))
	for _, c := range r.Changes {
		m[c.Key] = c.Action
	}
	return m
}
