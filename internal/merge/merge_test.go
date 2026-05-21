package merge

import (
	"testing"
)

func TestMerge_NoOverlap_CombinesBothMaps(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3"}

	res, err := Merge(a, b, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 3 {
		t.Errorf("expected 3 keys, got %d", len(res.Vars))
	}
	if len(res.Overrides) != 0 {
		t.Errorf("expected no overrides, got %v", res.Overrides)
	}
}

func TestMerge_PreferB_WinsOnConflict(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}

	res, err := Merge(a, b, Options{Strategy: PreferB})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["KEY"] != "from-b" {
		t.Errorf("expected 'from-b', got %q", res.Vars["KEY"])
	}
	if len(res.Overrides) != 1 || res.Overrides[0] != "KEY" {
		t.Errorf("expected override for KEY, got %v", res.Overrides)
	}
}

func TestMerge_PreferA_WinsOnConflict(t *testing.T) {
	a := map[string]string{"KEY": "from-a"}
	b := map[string]string{"KEY": "from-b"}

	res, err := Merge(a, b, Options{Strategy: PreferA})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Vars["KEY"] != "from-a" {
		t.Errorf("expected 'from-a', got %q", res.Vars["KEY"])
	}
	if len(res.Overrides) != 1 {
		t.Errorf("expected 1 override entry, got %v", res.Overrides)
	}
}

func TestMerge_ErrorOnConflict_ReturnsError(t *testing.T) {
	a := map[string]string{"KEY": "v1"}
	b := map[string]string{"KEY": "v2"}

	_, err := Merge(a, b, Options{Strategy: ErrorOnConflict})
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestMerge_IdenticalValues_NoOverride(t *testing.T) {
	a := map[string]string{"KEY": "same"}
	b := map[string]string{"KEY": "same"}

	res, err := Merge(a, b, Options{Strategy: ErrorOnConflict})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overrides) != 0 {
		t.Errorf("expected no overrides for identical values, got %v", res.Overrides)
	}
}

func TestMerge_EmptyMaps_ReturnsEmpty(t *testing.T) {
	res, err := Merge(map[string]string{}, map[string]string{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Vars) != 0 {
		t.Errorf("expected empty result, got %v", res.Vars)
	}
}
