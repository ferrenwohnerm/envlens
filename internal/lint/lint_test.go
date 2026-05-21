package lint_test

import (
	"testing"

	"github.com/yourorg/envlens/internal/lint"
)

func TestRun_NoFindings_WhenAllValid(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "abc123",
	}
	findings := lint.Run(env, lint.DefaultOptions())
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %+v", len(findings), findings)
	}
}

func TestRun_UppercaseKeys_DetectsLowercase(t *testing.T) {
	env := map[string]string{
		"database_url": "postgres://localhost/db",
		"API_KEY":      "abc123",
	}
	opts := lint.Options{UppercaseKeys: true}
	findings := lint.Run(env, opts)
	if !containsRule(findings, lint.RuleUppercaseKeys) {
		t.Error("expected uppercase-keys finding for 'database_url'")
	}
	if len(findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(findings))
	}
}

func TestRun_NoEmptyValues_DetectsBlank(t *testing.T) {
	env := map[string]string{
		"API_KEY": "",
		"HOST":    "   ",
		"PORT":    "8080",
	}
	opts := lint.Options{NoEmptyValues: true}
	findings := lint.Run(env, opts)
	count := countRule(findings, lint.RuleNoEmptyValues)
	if count != 2 {
		t.Errorf("expected 2 empty-value findings, got %d", count)
	}
}

func TestRun_NoWhitespaceKeys_DetectsSpaces(t *testing.T) {
	env := map[string]string{
		"MY KEY": "value",
		"VALID":  "ok",
	}
	opts := lint.Options{NoWhitespaceKeys: true}
	findings := lint.Run(env, opts)
	if !containsRule(findings, lint.RuleNoWhitespaceKeys) {
		t.Error("expected no-whitespace-keys finding")
	}
}

func TestRun_DisabledOptions_ProducesNoFindings(t *testing.T) {
	env := map[string]string{
		"lowercase_key": "",
	}
	opts := lint.Options{} // all rules disabled
	findings := lint.Run(env, opts)
	if len(findings) != 0 {
		t.Errorf("expected no findings with all rules disabled, got %d", len(findings))
	}
}

func TestRun_EmptyEnv_ReturnsNil(t *testing.T) {
	findings := lint.Run(map[string]string{}, lint.DefaultOptions())
	if len(findings) != 0 {
		t.Errorf("expected no findings for empty env, got %d", len(findings))
	}
}

func containsRule(findings []lint.Finding, rule lint.Rule) bool {
	for _, f := range findings {
		if f.Rule == rule {
			return true
		}
	}
	return false
}

func countRule(findings []lint.Finding, rule lint.Rule) int {
	count := 0
	for _, f := range findings {
		if f.Rule == rule {
			count++
		}
	}
	return count
}
