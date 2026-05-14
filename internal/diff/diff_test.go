package diff_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/diff"
)

func TestCompare_IdenticalMaps(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}

	r := diff.Compare(a, b)

	if r.HasDifferences() {
		t.Error("expected no differences for identical maps")
	}
	if len(r.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged keys, got %d", len(r.Unchanged))
	}
}

func TestCompare_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "val"}
	b := map[string]string{"FOO": "bar"}

	r := diff.Compare(a, b)

	if len(r.OnlyInA) != 1 {
		t.Errorf("expected 1 key only in A, got %d", len(r.OnlyInA))
	}
	if r.OnlyInA["ONLY_A"] != "val" {
		t.Errorf("unexpected value for ONLY_A: %q", r.OnlyInA["ONLY_A"])
	}
}

func TestCompare_OnlyInB(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "secret"}

	r := diff.Compare(a, b)

	if len(r.OnlyInB) != 1 {
		t.Errorf("expected 1 key only in B, got %d", len(r.OnlyInB))
	}
	if r.OnlyInB["ONLY_B"] != "secret" {
		t.Errorf("unexpected value for ONLY_B: %q", r.OnlyInB["ONLY_B"])
	}
}

func TestCompare_ChangedValues(t *testing.T) {
	a := map[string]string{"DB_HOST": "localhost", "PORT": "3000"}
	b := map[string]string{"DB_HOST": "prod.db.internal", "PORT": "3000"}

	r := diff.Compare(a, b)

	if len(r.Changed) != 1 {
		t.Errorf("expected 1 changed key, got %d", len(r.Changed))
	}
	pair, ok := r.Changed["DB_HOST"]
	if !ok {
		t.Fatal("expected DB_HOST to be in Changed")
	}
	if pair[0] != "localhost" || pair[1] != "prod.db.internal" {
		t.Errorf("unexpected changed values: %v", pair)
	}
	if len(r.Unchanged) != 1 {
		t.Errorf("expected 1 unchanged key, got %d", len(r.Unchanged))
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	r := diff.Compare(map[string]string{}, map[string]string{})

	if r.HasDifferences() {
		t.Error("expected no differences for two empty maps")
	}
}
