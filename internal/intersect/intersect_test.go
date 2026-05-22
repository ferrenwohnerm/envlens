package intersect_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/intersect"
)

func TestRun_NoOverlap_ReturnsEmptyResult(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3", "QUX": "4"}

	res := intersect.Run(a, b, intersect.DefaultOptions())

	if len(res.SharedKeys) != 0 {
		t.Errorf("expected 0 shared keys, got %d", len(res.SharedKeys))
	}
	if len(res.MatchingPairs) != 0 {
		t.Errorf("expected 0 matching pairs, got %d", len(res.MatchingPairs))
	}
	if len(res.DivergentKeys) != 0 {
		t.Errorf("expected 0 divergent keys, got %d", len(res.DivergentKeys))
	}
}

func TestRun_IdenticalMaps_AllKeysMatch(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "hello"}
	b := map[string]string{"FOO": "1", "BAR": "hello"}

	res := intersect.Run(a, b, intersect.DefaultOptions())

	if len(res.SharedKeys) != 2 {
		t.Errorf("expected 2 shared keys, got %d", len(res.SharedKeys))
	}
	if len(res.MatchingPairs) != 2 {
		t.Errorf("expected 2 matching pairs, got %d", len(res.MatchingPairs))
	}
	if len(res.DivergentKeys) != 0 {
		t.Errorf("expected 0 divergent keys, got %d", len(res.DivergentKeys))
	}
}

func TestRun_PartialOverlap_DivergentAndMatching(t *testing.T) {
	a := map[string]string{"FOO": "same", "BAR": "old", "ONLY_A": "x"}
	b := map[string]string{"FOO": "same", "BAR": "new", "ONLY_B": "y"}

	res := intersect.Run(a, b, intersect.DefaultOptions())

	if len(res.SharedKeys) != 2 {
		t.Errorf("expected 2 shared keys, got %d", len(res.SharedKeys))
	}
	if _, ok := res.MatchingPairs["FOO"]; !ok {
		t.Error("expected FOO in matching pairs")
	}
	if len(res.DivergentKeys) != 1 || res.DivergentKeys[0] != "BAR" {
		t.Errorf("expected [BAR] as divergent keys, got %v", res.DivergentKeys)
	}
}

func TestRun_CaseFold_MatchesCaseInsensitiveKeys(t *testing.T) {
	a := map[string]string{"foo": "val"}
	b := map[string]string{"FOO": "val"}

	opts := intersect.Options{CaseFold: true}
	res := intersect.Run(a, b, opts)

	if len(res.SharedKeys) != 1 {
		t.Errorf("expected 1 shared key with CaseFold, got %d", len(res.SharedKeys))
	}
}

func TestRun_CaseFoldDisabled_NoMatch(t *testing.T) {
	a := map[string]string{"foo": "val"}
	b := map[string]string{"FOO": "val"}

	res := intersect.Run(a, b, intersect.DefaultOptions())

	if len(res.SharedKeys) != 0 {
		t.Errorf("expected 0 shared keys without CaseFold, got %d", len(res.SharedKeys))
	}
}

func TestRun_EmptyMaps_ReturnsEmptyResult(t *testing.T) {
	res := intersect.Run(map[string]string{}, map[string]string{}, intersect.DefaultOptions())

	if len(res.SharedKeys) != 0 || len(res.MatchingPairs) != 0 || len(res.DivergentKeys) != 0 {
		t.Error("expected fully empty result for empty input maps")
	}
}
