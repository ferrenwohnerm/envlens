package filter_test

import (
	"testing"

	"github.com/user/envlens/internal/diff"
	"github.com/user/envlens/internal/filter"
)

func makeResults() map[string]diff.Result {
	return map[string]diff.Result{
		"APP_NAME":    {Status: diff.StatusEqual, ValueA: "envlens", ValueB: "envlens"},
		"APP_ENV":     {Status: diff.StatusChanged, ValueA: "staging", ValueB: "production"},
		"DB_HOST":     {Status: diff.StatusOnlyInA, ValueA: "localhost", ValueB: ""},
		"SECRET_KEY":  {Status: diff.StatusOnlyInB, ValueA: "", ValueB: "s3cr3t"},
		"LOG_LEVEL":   {Status: diff.StatusEqual, ValueA: "info", ValueB: "info"},
	}
}

func TestApply_NoOptions_ReturnsAll(t *testing.T) {
	results := makeResults()
	out := filter.Apply(results, filter.Options{})
	if len(out) != len(results) {
		t.Errorf("expected %d entries, got %d", len(results), len(out))
	}
}

func TestApply_PrefixFilter(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{Prefix: "APP_"})
	if len(out) != 2 {
		t.Errorf("expected 2 entries with prefix APP_, got %d", len(out))
	}
	for k := range out {
		if k != "APP_NAME" && k != "APP_ENV" {
			t.Errorf("unexpected key %q in prefix-filtered results", k)
		}
	}
}

func TestApply_OnlyChanged(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{OnlyChanged: true})
	if len(out) != 1 {
		t.Errorf("expected 1 changed entry, got %d", len(out))
	}
	if _, ok := out["APP_ENV"]; !ok {
		t.Error("expected APP_ENV in changed results")
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{OnlyMissing: true})
	if len(out) != 2 {
		t.Errorf("expected 2 missing entries, got %d", len(out))
	}
}

func TestApply_ExcludeKeys(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{ExcludeKeys: []string{"SECRET_KEY", "DB_HOST"}})
	if _, found := out["SECRET_KEY"]; found {
		t.Error("SECRET_KEY should have been excluded")
	}
	if _, found := out["DB_HOST"]; found {
		t.Error("DB_HOST should have been excluded")
	}
	if len(out) != 3 {
		t.Errorf("expected 3 entries after exclusion, got %d", len(out))
	}
}

func TestApply_PrefixCaseInsensitive(t *testing.T) {
	out := filter.Apply(makeResults(), filter.Options{Prefix: "app_"})
	if len(out) != 2 {
		t.Errorf("expected 2 entries with case-insensitive prefix, got %d", len(out))
	}
}
