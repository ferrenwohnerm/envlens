package dedup

import (
	"testing"
)

func TestRun_NoOverlap_CombinesAllMaps(t *testing.T) {
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	c := map[string]string{"C": "3"}

	res, err := Run([]map[string]string{a, b, c}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 3 {
		t.Errorf("expected 3 vars, got %d", len(res.Vars))
	}
	if len(res.Duplicates) != 0 {
		t.Errorf("expected no duplicates, got %v", res.Duplicates)
	}
}

func TestRun_PreferFirst_KeepsFirstValue(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}

	opts := DefaultOptions()
	opts.PreferFirst = true

	res, err := Run([]map[string]string{a, b}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["KEY"] != "from-a" {
		t.Errorf("expected 'from-a', got %q", res.Vars["KEY"])
	}
	if len(res.Duplicates) != 1 || res.Duplicates[0] != "KEY" {
		t.Errorf("expected [KEY] in duplicates, got %v", res.Duplicates)
	}
}

func TestRun_PreferLast_KeepsLastValue(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}

	opts := Options{PreferFirst: false}

	res, err := Run([]map[string]string{a, b}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["KEY"] != "from-b" {
		t.Errorf("expected 'from-b', got %q", res.Vars["KEY"])
	}
}

func TestRun_ErrorOnDuplicate_ReturnsError(t *testing.T) {
	a := map[string]string{"KEY": "1"}
	b := map[string]string{"KEY": "2"}

	opts := Options{ErrorOnDuplicate: true}

	_, err := Run([]map[string]string{a, b}, opts)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	dke, ok := err.(*DuplicateKeyError)
	if !ok {
		t.Fatalf("expected *DuplicateKeyError, got %T", err)
	}
	if dke.Key != "KEY" {
		t.Errorf("expected key 'KEY', got %q", dke.Key)
	}
}

func TestRun_EmptyInput_ReturnsEmptyResult(t *testing.T) {
	res, err := Run([]map[string]string{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 0 {
		t.Errorf("expected empty vars, got %v", res.Vars)
	}
}

func TestRun_DuplicatesAreSorted(t *testing.T) {
	a := map[string]string{"Z": "1", "A": "1", "M": "1"}
	b := map[string]string{"Z": "2", "A": "2", "M": "2"}

	res, err := Run([]map[string]string{a, b}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"A", "M", "Z"}
	for i, k := range res.Duplicates {
		if k != expected[i] {
			t.Errorf("pos %d: expected %q, got %q", i, expected[i], k)
		}
	}
}
